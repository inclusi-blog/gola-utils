package model

type IdToken struct {
	UserId          string `json:"userId"`
	Username        string `json:"username"`
	Email           string `json:"email"`
	Subject         string `json:"subject"`
	AccessTokenHash string `json:"at_hash"`
}


