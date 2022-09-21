package handlers

import (
	"fmt"
	utils2 "github.com/aabri-ankorstore/cli-auth/pkg/utils"
	"net/http"
)

func (h *Auth) LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Cache-Control", "no-cache") // See https://github.com/okta/samples-golang/issues/20

	utils2.Nonce, _ = utils2.GenerateNonce()
	var url string

	q := r.URL.Query()
	q.Add("client_id", utils2.ClientID)
	q.Add("response_type", "code")
	q.Add("response_mode", "query")
	q.Add("scope", "openid profile email")
	q.Add("redirect_uri", utils2.RedirectUri)
	q.Add("state", utils2.State)
	q.Add("nonce", utils2.Nonce)
	url = fmt.Sprintf("%s/v1/authorize?%s", utils2.ISSUER, q.Encode())
	http.Redirect(w, r, url, http.StatusFound)
}
