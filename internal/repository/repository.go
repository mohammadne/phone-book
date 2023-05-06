package repository

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"strings"

	"github.com/MohammadNE/PhoneBook/internal/models"
	"github.com/MohammadNE/PhoneBook/pkg/rdbms"
	"github.com/MohammadNE/PhoneBook/pkg/utils"
	"go.uber.org/zap"
)

type Repository interface {
	Migrate(models.Migrate) error

	CreateUser(ctx context.Context, user *models.User) error
	FindUserByEmail(ctx context.Context, email string) (*models.User, error)
	FindUserByEmailAndPassword(ctx context.Context, email, password string) (*models.User, error)

	DeleteContact(ctx context.Context, userId, contactId uint64) error
}

type repository struct {
	logger *zap.Logger
	rdbms  rdbms.RDBMS
}

func New(logger *zap.Logger, rdbms rdbms.RDBMS) Repository {
	r := &repository{logger: logger, rdbms: rdbms}

	return r
}

//go:embed migrations
var migrations embed.FS

func (r *repository) Migrate(direction models.Migrate) error {
	files, err := fs.ReadDir(migrations, "migrations")
	if err != nil {
		return fmt.Errorf("error reading migrations directory:\n%v", err)
	}

	result := make([]string, 0, len(files)/2)

	for _, file := range files {
		splits := strings.Split(file.Name(), ".")
		if splits[1] == string(direction) {
			result = append(result, file.Name())
		}
	}

	result = utils.Sort(result)

	for index := 0; index < len(result); index++ {
		file := "migrations/"

		if direction == models.Up {
			file += result[index]
		} else {
			file += result[len(result)-index-1]
		}

		data, err := fs.ReadFile(migrations, file)
		if err != nil {
			return fmt.Errorf("error reading migration file: %s\n%v", file, err)
		}

		if err := r.rdbms.Execute(string(data), []any{}); err != nil {
			return fmt.Errorf("error migrating the file: %s\n%v", file, err)
		}
	}

	return nil
}
