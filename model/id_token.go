package model

type IdToken struct {
	UserId          string `json:"id"`
	Username        string `json:"username"`
	Email           string `json:"email"`
	Subject         string `json:"subject"`
	AccessTokenHash string `json:"at_hash"`
}
