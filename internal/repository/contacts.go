package repository

import (
	"strconv"
	"strings"

	"github.com/lib/pq"
	"github.com/mohammadne/phone-book/internal/models"
	"github.com/mohammadne/phone-book/pkg/crypto"
	"go.uber.org/zap"
)

const QueryCreateContact = `
INSERT INTO contacts(name, phones, description, user_id) VALUES($1, $2, $3, $4) 
RETURNING id;`

func (r *repository) CreateContact(userId uint64, contact *models.Contact) error {
	in := []interface{}{contact.Name, pq.Array(contact.Phones), contact.Description, userId}
	out := []any{&contact.Id}
	if err := r.rdbms.QueryRow(QueryCreateContact, in, out); err != nil {
		r.logger.Error("Error inserting contact", zap.Error(err))
		return err
	}
	return nil
}

const QueryGetContactById = `
SELECT name, phones, description
FROM contacts
WHERE user_id=$1 AND id=$2;`

func (r *repository) GetContactById(userId, contactId uint64) (*models.Contact, error) {
	contact := models.Contact{Id: contactId}

	in := []any{userId, contactId}
	out := []any{&contact.Name, pq.Array(&contact.Phones), &contact.Description}
	if err := r.rdbms.QueryRow(QueryGetContactById, in, out); err != nil {
		r.logger.Error("Error get contact by id", zap.Error(err))
		return nil, err
	}

	return &contact, nil
}

const QueryUpdateContact = `
UPDATE contacts 
SET name=$1, phones=$2, description=$3 
WHERE user_id=$4 AND id=$5;`

func (r *repository) UpdateContact(userId uint64, contact *models.Contact) error {
	in := []any{contact.Name, pq.Array(contact.Phones), contact.Description, userId, contact.Id}
	if err := r.rdbms.Execute(QueryUpdateContact, in); err != nil {
		r.logger.Error("Error updating contact", zap.Error(err))
		return err
	}
	return nil
}

const QueryDeleteContact = `
DELETE FROM contacts 
WHERE user_id=$1 AND id=$2;`

func (r *repository) DeleteContact(userId, contactId uint64) error {
	in := []interface{}{userId, contactId}
	if err := r.rdbms.Execute(QueryDeleteContact, in); err != nil {
		r.logger.Error("Error deleting contact", zap.Uint64("user-id", userId), zap.Uint64("contact-id", contactId), zap.Error(err))
		return err
	}
	return nil
}

const QueryGetContacts = `
SELECT id, name, phones, description
FROM contacts
WHERE 
	user_id=$1 AND 
	id > $2 AND 
	name LIKE '%' || $3 || '%'
ORDER BY id
FETCH NEXT $4 ROWS ONLY;`

func (r *repository) GetContacts(userId uint64, encryptedCursor, search string, limit int) ([]models.Contact, string, error) {
	var id uint64 = 0

	if limit < r.config.Limit.Min {
		limit = r.config.Limit.Min
	} else if limit > r.config.Limit.Max {
		limit = r.config.Limit.Max
	}

	// decrypt cursor
	if len(encryptedCursor) != 0 {
		cursor, err := crypto.Decrypt(encryptedCursor, r.config.CursorSecret)
		if err != nil {
			panic(err)
		}

		splits := strings.Split(cursor, ",")
		if len(splits) != 1 {
			panic("err")
		}

		id, err = strconv.ParseUint(splits[0], 10, 64)
		if err != nil {
			panic(err)
		}
	}

	contacts := make([]models.Contact, limit)
	out := make([][]any, limit)

	for index := 0; index < limit; index++ {
		out[index] = []any{&contacts[index].Id, &contacts[index].Name, pq.Array(&contacts[index].Phones), &contacts[index].Description}
	}

	in := []any{userId, id, search, limit}
	if err := r.rdbms.Query(QueryGetContacts, in, out); err != nil {
		r.logger.Error("Error query contacts", zap.Error(err))
		return nil, "", err
	}

	if len(contacts) == 0 {
		return contacts, "", nil
	}

	var lastContact models.Contact

	for index := limit - 1; index >= 0; index-- {
		if contacts[index].Id != 0 {
			lastContact = contacts[index]
			break
		} else {
			contacts = contacts[:index]
		}
	}

	if lastContact.Id == 0 {
		return contacts, "", nil
	}

	cursor := strconv.FormatUint(lastContact.Id, 10)

	// encrypt cursor
	encryptedCursor, err := crypto.Encrypt(cursor, r.config.CursorSecret)
	if err != nil {
		panic(err)
	}

	return contacts, encryptedCursor, nil
}
