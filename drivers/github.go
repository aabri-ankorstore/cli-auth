package drivers

import (
	"bytes"
	"encoding/json"
	"fmt"
	utils2 "github.com/aabri-ankorstore/cli-auth/utils"
	"github.com/ankorstore/ankorstore-cli-core/pkg/util"
	"github.com/go-errors/errors"
	"github.com/michaeljs1990/sqlitestore"
	verifier "github.com/okta/okta-jwt-verifier-golang"
	"github.com/rs/zerolog/log"
	"github.com/skratchdot/open-golang/open"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type Github struct {
	SessionStore     *sqlitestore.SqliteStore
	SessionStoreName string
	Nonce            string
	State            string
}

func NewGithubClient() Manager {
	utils2.Nonce, _ = utils2.GenerateNonce()
	utils2.State = utils2.GenerateState()
	return &Github{
		SessionStore:     utils2.SessionStore,
		SessionStoreName: utils2.CookieName,
		Nonce:            utils2.Nonce,
		State:            utils2.State,
	}
}

func (g *Github) InformUserAndOpenBrowser() error {
	log.Info().Msg("Opening browser for code entry...")
	// Wait a few seconds to give user a chance to check out the printed user code.
	time.Sleep(3 * time.Second)
	var url string
	r, _ := http.NewRequest(http.MethodGet, host, nil)
	q := r.URL.Query()
	q.Add("client_id", utils2.ClientID)
	q.Add("redirect_uri", utils2.RedirectUri)
	q.Add("state", g.State)
	q.Add("nonce", g.Nonce)
	url = fmt.Sprintf("%s/authorize?%s", utils2.ISSUER, q.Encode())
	err := open.Run(url)
	if err != nil {
		return err
	}
	return nil
}
func (g *Github) ExchangeCode(w http.ResponseWriter, r *http.Request) (Exchange, error) {
	// Check the state that was returned to the query string is the same as the above state
	if r.URL.Query().Get("state") != utils2.State {
		fmt.Fprintln(w, "The state was not as expected")
		return Exchange{}, nil
	}
	// Make sure the code was provided
	if r.URL.Query().Get("code") == "" {
		fmt.Fprintln(w, "The code was not returned or is not accessible")
		return Exchange{}, nil
	}
	code := r.URL.Query().Get("code")

	requestBodyMap := map[string]string{"client_id": utils2.ClientID, "client_secret": utils2.ClientSECRET, "code": code}
	requestJSON, _ := json.Marshal(requestBodyMap)

	url := fmt.Sprintf("%s/access_token", utils2.ISSUER)
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

	session, err := utils2.SessionStore.Get(r, utils2.CookieName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	//session.Values["id_token"] = exchange.IdToken
	session.Values["access_token"] = exchange.AccessToken
	_ = session.Save(r, w)
	return exchange, nil
}
func (g *Github) VerifyToken(t string) (*verifier.Jwt, error) {
	tv := map[string]string{}
	tv["nonce"] = utils2.Nonce
	tv["aud"] = utils2.ClientID
	jv := verifier.JwtVerifier{
		Issuer:           utils2.ISSUER,
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
	session, err := utils2.SessionStore.Get(r, utils2.CookieName)
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
	var err error
	dirs := util.NewDirs()
	authDir := fmt.Sprintf("%s/%s", dirs.GetPluginsDir(), utils2.PluginPath)
	if _, err := os.Stat(authDir); os.IsNotExist(err) {
		authDir = "."
	}
	utils2.SessionStore, err = sqlitestore.NewSqliteStore(
		fmt.Sprintf("%s/auth", authDir),
		"sessions",
		"/",
		3600,
		[]byte(utils2.CookieName))
	if err != nil {
		panic(err)
	}
}
