package entities

// AccessToken ...
type AccessToken struct {
	ClientID    string `json:"client_id"`
	AccessToken string `json:"access_token"`
	IdToken     string `json:"id_token"`
}
