package models

type User struct {
	Id        uint64    `json:"Id"`
	Email     string    `json:"email"`
	Password  string    `json:"password,omitempty"`
	Contacts  []Contact `json:"contacts,omitempty"`
	CreatedAt string    `json:"created_at"`
}

func (c User) Marshal() *User {
	return &User{
		Id:        c.Id,
		Email:     c.Email,
		Contacts:  c.Contacts,
		CreatedAt: c.CreatedAt,
	}
}
