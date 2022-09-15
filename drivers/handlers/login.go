package handlers

import (
	"fmt"
	"github.com/aabri-ankorstore/cli-auth/utils"
	"net/http"
)

func (h *Auth) LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Cache-Control", "no-cache") // See https://github.com/okta/samples-golang/issues/20

	utils.Nonce, _ = utils.GenerateNonce()
	var url string

	q := r.URL.Query()
	q.Add("client_id", utils.ClientID)
	q.Add("response_type", "code")
	q.Add("response_mode", "query")
	q.Add("scope", "openid profile email")
	q.Add("redirect_uri", utils.RedirectUri)
	q.Add("state", utils.State)
	q.Add("nonce", utils.Nonce)
	url = fmt.Sprintf("%s/v1/authorize?%s", utils.ISSUER, q.Encode())
	http.Redirect(w, r, url, http.StatusFound)
}
