package repository

import (
	"embed"
	"fmt"
	"io/fs"
	"strings"

	"github.com/mohammadne/phone-book/internal/models"
	"github.com/mohammadne/phone-book/pkg/rdbms"
	"github.com/mohammadne/phone-book/pkg/utils"
	"go.uber.org/zap"
)

type Repository interface {
	Migrate(models.Migrate) error

	CreateUser(user *models.User) error
	GetUserByEmail(email string) (*models.User, error)
	GetUserByEmailAndPassword(email, password string) (*models.User, error)

	CreateContact(userId uint64, contect *models.Contact) error
	GetContactById(userId, contactId uint64) (*models.Contact, error)
	UpdateContact(userId uint64, contact *models.Contact) error
	DeleteContact(userId, contactId uint64) error
	GetContacts(userId uint64, encryptedCursor, search string, limit int) ([]models.Contact, string, error)
}

type repository struct {
	logger *zap.Logger
	config *Config
	rdbms  rdbms.RDBMS
}

func New(logger *zap.Logger, cfg *Config, rdbms rdbms.RDBMS) Repository {
	r := &repository{logger: logger, config: cfg, rdbms: rdbms}

	return r
}

//go:embed migrations
var migrations embed.FS

func (r *repository) Migrate(direction models.Migrate) error {
	files, err := fs.ReadDir(migrations, "migrations")
	if err != nil {
		return fmt.Errorf("Error reading migrations directory:\n%v", err)
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
			return fmt.Errorf("Error reading migration file: %s\n%v", file, err)
		}

		if err := r.rdbms.Execute(string(data), []any{}); err != nil {
			return fmt.Errorf("Error migrating the file: %s\n%v", file, err)
		}
	}

	return nil
}
