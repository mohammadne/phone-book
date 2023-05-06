package repository

import (
	"context"

	"go.uber.org/zap"
)

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
