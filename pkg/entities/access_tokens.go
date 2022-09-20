package entities

// AccessToken ...
type AccessToken struct {
	AccountID   string `json:"account_id"`
	AccessToken string `json:"access_token"`
	TokenExpiry string `json:"token_expiry"`
	RaptToken   string `json:"rapt_token"`
	IdToken     string `json:"id_token"`
}
