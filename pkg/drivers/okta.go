package drivers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aabri-ankorstore/cli-auth/pkg/utils"
	"github.com/gorilla/sessions"
	verifier "github.com/okta/okta-jwt-verifier-golang"
	"github.com/rs/zerolog/log"
	"github.com/skratchdot/open-golang/open"
	"io/ioutil"
	"net/http"
	"time"
)

type Okta struct {
	SessionStore     *sessions.CookieStore
	SessionStoreName string
	Nonce            string
	State            string
}

func NewOktaClient() Manager {
	utils.Nonce, _ = utils.GenerateNonce()
	utils.State = utils.GenerateState()
	return &Okta{
		SessionStore:     utils.SessionStore,
		SessionStoreName: utils.CookieName,
		Nonce:            utils.Nonce,
		State:            utils.State,
	}
}
func (g *Okta) InformUserAndOpenBrowser() error {
	log.Info().Msg("Opening browser for code entry...")
	// Wait a few seconds to give user a chance to check out the printed user code.
	time.Sleep(3 * time.Second)
	var url string
	r, _ := http.NewRequest(http.MethodGet, host, nil)
	q := r.URL.Query()
	q.Add("client_id", utils.ClientID)
	q.Add("response_type", "code")
	q.Add("response_mode", "query")
	q.Add("scope", "openid profile email")
	q.Add("redirect_uri", utils.RedirectUri)
	q.Add("state", g.State)
	q.Add("nonce", g.Nonce)
	url = fmt.Sprintf("%s/v1/authorize?%s", utils.ISSUER, q.Encode())
	err := open.Run(url)
	if err != nil {
		return err
	}
	return nil
}
func (g *Okta) ExchangeCode(w http.ResponseWriter, r *http.Request) (Exchange, error) {
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

	authHeader := base64.StdEncoding.EncodeToString(
		[]byte(fmt.Sprintf("%s:%s", utils.ClientID, utils.ClientSECRET)))

	q := r.URL.Query()
	q.Add("grant_type", "authorization_code")
	q.Set("code", code)
	q.Add("redirect_uri", utils.RedirectUri)

	url := fmt.Sprintf("%s/v1/token?%s", utils.ISSUER, q.Encode())

	req, _ := http.NewRequest("POST", url, bytes.NewReader([]byte("")))
	h := req.Header
	h.Add("Authorization", fmt.Sprintf("Basic %s", authHeader))
	h.Add("Accept", "application/json")
	h.Add("Content-Type", "application/x-www-form-urlencoded")
	h.Add("Connection", "close")
	h.Add("Content-Length", "0")

	client := &http.Client{}
	resp, _ := client.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	var exchange Exchange
	_ = json.Unmarshal(body, &exchange)

	// Verification
	_, verificationError := g.VerifyToken(exchange.IdToken)
	if verificationError != nil {
		fmt.Println(verificationError)
	}

	if verificationError == nil {
		session, err := utils.SessionStore.Get(r, utils.CookieName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		session.Values["id_token"] = exchange.IdToken
		session.Values["access_token"] = exchange.AccessToken
		err = session.Save(r, w)
		if err != nil {
			return Exchange{}, err
		}
	}

	return exchange, nil
}
func (g *Okta) VerifyToken(t string) (*verifier.Jwt, error) {
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
func (g *Okta) GetProfile(r *http.Request) (map[string]string, error) {
	m := make(map[string]string)
	session, err := utils.SessionStore.Get(r, utils.CookieName)
	if err != nil || session.Values["access_token"] == nil || session.Values["access_token"] == "" {
		return m, errors.New("please provide a valid access token")
	}

	token := session.Values["access_token"].(string)
	reqUrl := fmt.Sprintf("%s/v1/userinfo", utils.ISSUER)
	req, _ := http.NewRequest("GET", reqUrl, bytes.NewReader([]byte("")))
	h := req.Header
	h.Add("Authorization", fmt.Sprintf("Bearer %s", token))
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
