package repository

import (
	"errors"

	"github.com/mohammadne/phone-book/internal/models"
	"github.com/mohammadne/phone-book/pkg/rdbms"
	"go.uber.org/zap"
)

const QueryCreateUser = `
INSERT INTO users(email, password) VALUES($1, $2) 
RETURNING id;`

func (r *repository) CreateUser(user *models.User) error {
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

const QueryGetUserDetailsById = `
SELECT users.id, users.password, users.created_at, contacts.name, contacts.phones, contacts.description
FROM users
JOIN contacts ON users.id = contacts.user_id
WHERE 
	users.id=$1 AND 
	contacts.id > $2 AND 
	contacts.name LIKE '%' || $3 || '%'
ORDER BY contacts.id
FETCH NEXT $4 ROWS ONLY;`

const QueryGetUserByEmail = `
SELECT id, password, created_at
FROM users
WHERE email=$1;`

func (r *repository) GetUserByEmail(email string) (*models.User, error) {
	user := &models.User{Email: email}

	in := []interface{}{email}
	out := []interface{}{&user.Id, &user.Password, &user.CreatedAt}
	if err := r.rdbms.QueryRow(QueryGetUserByEmail, in, out); err != nil {
		if err.Error() == rdbms.ErrNotFound {
			return nil, err
		}

		r.logger.Error("Error find user by email", zap.Error(err))
		return nil, err
	}

	return user, nil
}

const QueryGetUserByEmailAndPassword = `
SELECT id, created_at 
FROM users 
WHERE email=$1 AND password=$2;`

func (r *repository) GetUserByEmailAndPassword(email, password string) (*models.User, error) {
	user := &models.User{Email: email, Password: password}

	in := []interface{}{email, password}
	out := []interface{}{&user.Id, &user.CreatedAt}
	if err := r.rdbms.QueryRow(QueryGetUserByEmailAndPassword, in, out); err != nil {
		r.logger.Error("Error find user by email and password", zap.Error(err))
		return nil, err
	}

	return user, nil
}
