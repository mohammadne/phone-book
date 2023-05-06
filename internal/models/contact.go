package models

type Contact struct {
	Id          uint64   `json:"id"`
	Name        string   `json:"name"`
	Phones      []string `json:"phones"`
	Description string   `json:"description,omitempty"`
}

func (c *Contact) IsValid() bool {
	if len(c.Name) == 0 || len(c.Phones) == 0 {
		return false
	}
	return true
}

func (newContact *Contact) Update(oldContact *Contact) {
	newContact.Id = oldContact.Id

	if len(newContact.Name) == 0 {
		newContact.Name = oldContact.Name
	}

	if len(newContact.Phones) == 0 {
		newContact.Phones = oldContact.Phones
	}

	if len(newContact.Description) == 0 {
		newContact.Description = oldContact.Description
	}
}
