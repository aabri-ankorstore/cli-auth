package drivers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/aabri-ankorstore/cli-auth/pkg/utils"
	"github.com/go-errors/errors"
	"github.com/gorilla/sessions"
	verifier "github.com/okta/okta-jwt-verifier-golang"
	"github.com/skratchdot/open-golang/open"
	"io/ioutil"
	"net/http"
	"time"
)

type Github struct {
	SessionStore     *sessions.CookieStore
	SessionStoreName string
	Nonce            string
	State            string
}

func NewGithubClient() Manager {
	utils.Nonce, _ = utils.GenerateNonce()
	utils.State = utils.GenerateState()
	return &Github{
		SessionStore:     utils.SessionStore,
		SessionStoreName: utils.CookieName,
		Nonce:            utils.Nonce,
		State:            utils.State,
	}
}

func (g *Github) InformUserAndOpenBrowser() error {
	fmt.Println("Opening browser for code entry...")
	// Wait a few seconds to give user a chance to check out the printed user code.
	time.Sleep(3 * time.Second)
	var url string
	r, _ := http.NewRequest(http.MethodGet, host, nil)
	q := r.URL.Query()
	q.Add("client_id", utils.ClientID)
	q.Add("redirect_uri", utils.RedirectUri)
	q.Add("state", g.State)
	q.Add("nonce", g.Nonce)
	url = fmt.Sprintf("%s/authorize?%s", utils.ISSUER, q.Encode())
	err := open.Run(url)
	if err != nil {
		return err
	}
	return nil
}
func (g *Github) ExchangeCode(w http.ResponseWriter, r *http.Request) (Exchange, error) {
	// Check the state that was returned to the query string is the same as the above state
	if r.URL.Query().Get("state") != utils.State {
		fmt.Fprintln(w, "The state was not as expected")
		return Exchange{}, nil
	}
	// Make sure the code was provided
	if r.URL.Query().Get("code") == "" {
		fmt.Fprintln(w, "The code was not returned or is not accessible")
		return Exchange{}, nil
	}
	code := r.URL.Query().Get("code")

	requestBodyMap := map[string]string{"client_id": utils.ClientID, "client_secret": utils.ClientSECRET, "code": code}
	requestJSON, _ := json.Marshal(requestBodyMap)

	url := fmt.Sprintf("%s/access_token", utils.ISSUER)
	req, reqerr := http.NewRequest("POST", url, bytes.NewBuffer(requestJSON))
	if reqerr != nil {
		return Exchange{}, reqerr
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, _ := client.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	var exchange Exchange
	_ = json.Unmarshal(body, &exchange)

	session, err := utils.SessionStore.Get(r, utils.CookieName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	//session.Values["id_token"] = exchange.IdToken
	session.Values["access_token"] = exchange.AccessToken
	err = session.Save(r, w)
	if err != nil {
		return Exchange{}, err
	}
	return exchange, nil
}
func (g *Github) VerifyToken(t string) (*verifier.Jwt, error) {
	tv := map[string]string{}
	tv["nonce"] = utils.Nonce
	tv["aud"] = utils.ClientID
	jv := verifier.JwtVerifier{
		Issuer:           utils.ISSUER,
		ClaimsToValidate: tv,
	}

	result, err := jv.New().VerifyIdToken(t)
	if err != nil {
		return nil, fmt.Errorf("%s", err)
	}

	if result != nil {
		return result, nil
	}

	return nil, fmt.Errorf("token could not be verified: %s", "")
}
func (g *Github) GetProfile(r *http.Request) (map[string]string, error) {
	m := make(map[string]string)
	session, err := utils.SessionStore.Get(r, utils.CookieName)
	if err != nil || session.Values["access_token"] == nil || session.Values["access_token"] == "" {
		return m, errors.New("Please provide a valid access token")
	}
	token := session.Values["access_token"].(string)
	reqUrl := "https://api.github.com/user"
	req, reqerr := http.NewRequest("GET", reqUrl, nil)
	if reqerr != nil {
		return nil, reqerr
	}
	h := req.Header
	h.Add("Authorization", fmt.Sprintf("token %s", token))
	h.Add("Accept", "application/json")

	client := &http.Client{}
	resp, _ := client.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	_ = json.Unmarshal(body, &m)
	return m, nil
}

func init() {
	utils.SessionStore = sessions.NewCookieStore([]byte(utils.CookieName))
	utils.SessionStore.Options = &sessions.Options{
		Path:     "/",      // to match all requests
		MaxAge:   3600 * 1, // 1 hour
		HttpOnly: true,
	}
}
