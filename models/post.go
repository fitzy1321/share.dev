package models

type Post struct {
	ID      string `json:"id"`
	To      string `json:"to"`
	From    string `json:"from"`
	Message string `json:"Message"`
}
