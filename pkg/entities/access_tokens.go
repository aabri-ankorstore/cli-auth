package entities

// AccessToken ...
type AccessToken struct {
	AccountID   string `json:"account_id"`
	AccessToken string `json:"access_token"`
	IdToken     string `json:"id_token"`
}
