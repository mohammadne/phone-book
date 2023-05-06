package models

type Contact struct {
	Id          uint64   `json:"id"`
	Name        string   `json:"name"`
	Phones      []string `json:"phones"`
	Description string   `json:"description,omitempty"`
}
