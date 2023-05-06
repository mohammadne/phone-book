package repository

import (
	"context"

	"go.uber.org/zap"
)

const QueryDeleteContact = `
DELETE FROM contacts 
WHERE id=$1;`

func (r *repository) DeleteContact(ctx context.Context, id uint64) error {
	in := []interface{}{id}
	if err := r.rdbms.Execute(QueryDeleteContact, in); err != nil {
		r.logger.Error("Error deleting contact", zap.Uint64("id", id), zap.Error(err))
		return err
	}
	return nil
}
