package entity

type Account struct {
	Id      string    `json:"id,omitempty"`
	Account string `json:"account"`
	Pwd     string `json:"pwd"`
	Created int64  `json:"created,omitempty"`
}

