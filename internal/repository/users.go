package repository

import (
	"context"
	"errors"

	"github.com/MohammadNE/PhoneBook/internal/models"
	"github.com/MohammadNE/PhoneBook/pkg/rdbms"
	"go.uber.org/zap"
)

const QueryCreateUser = `
INSERT INTO users(email, password) VALUES($1, $2) 
RETURNING id;`

func (r *repository) CreateUser(ctx context.Context, user *models.User) error {
	if len(user.Email) == 0 || len(user.Password) == 0 {
		return errors.New("Insufficient information for user")
	}

	in := []any{user.Email, user.Password}
	out := []any{&user.Id}
	if err := r.rdbms.QueryRow(QueryCreateUser, in, out); err != nil {
		r.logger.Error("Error inserting author", zap.Error(err))
		return err
	}

	return nil
}

const QueryFindUserByEmail = `
SELECT id, password, created_at
FROM users
WHERE email=$1;`

func (r *repository) FindUserByEmail(ctx context.Context, email string) (*models.User, error) {
	user := &models.User{Email: email}

	args := []interface{}{email}
	dest := []interface{}{&user.Id, &user.Password, &user.CreatedAt}
	if err := r.rdbms.QueryRow(QueryFindUserByEmail, args, dest); err != nil {
		if err.Error() == rdbms.ErrNotFound {
			return nil, err
		}

		r.logger.Error("Error find user by email", zap.Error(err))
		return nil, err
	}

	return user, nil
}

const QueryFindUserByEmailAndPassword = `
SELECT id, created_at 
FROM users 
WHERE email=$1 AND password=$2;`

func (r *repository) FindUserByEmailAndPassword(ctx context.Context, email, password string) (*models.User, error) {
	user := &models.User{Email: email, Password: password}

	args := []interface{}{email, password}
	dest := []interface{}{&user.Id, &user.CreatedAt}
	if err := r.rdbms.QueryRow(QueryFindUserByEmailAndPassword, args, dest); err != nil {
		r.logger.Error("Error find user by email and password", zap.Error(err))
		return nil, err
	}

	return user, nil
}
