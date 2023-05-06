package repository

import (
	"context"

	"github.com/MohammadNE/PhoneBook/internal/models"
	"go.uber.org/zap"
)

const QueryCreateContact = `
INSERT INTO contacts(name, phones, description, user_id) VALUES($1, $2, $3, $4) 
RETURNING id;`

func (r *repository) CreateContact(ctx context.Context, userId uint64, contact *models.Contact) error {
	in := []interface{}{contact.Name, contact.Phones, contact.Description, userId}
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

func (r *repository) GetContactById(ctx context.Context, userId, contactId uint64) (*models.Contact, error) {
	contact := models.Contact{Id: contactId}

	in := []any{userId, contactId}
	out := []any{&contact.Name, &contact.Phones, &contact.Description}
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

func (r *repository) UpdateContact(ctx context.Context, userId uint64, contact *models.Contact) error {
	in := []interface{}{contact.Name, contact.Phones, contact.Description, userId, contact.Id}
	if err := r.rdbms.Execute(QueryCreateContact, in); err != nil {
		r.logger.Error("Error updating contact", zap.Error(err))
		return err
	}
	return nil
}

const QueryDeleteContact = `
DELETE FROM contacts 
WHERE user_id=$1 AND id=$2;`

func (r *repository) DeleteContact(ctx context.Context, userId, contactId uint64) error {
	in := []interface{}{userId, contactId}
	if err := r.rdbms.Execute(QueryDeleteContact, in); err != nil {
		r.logger.Error("Error deleting contact", zap.Uint64("user-id", userId), zap.Uint64("contact-id", contactId), zap.Error(err))
		return err
	}
	return nil
}
