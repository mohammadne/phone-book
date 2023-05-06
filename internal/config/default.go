package config

import (
	"time"

	"github.com/MohammadNE/PhoneBook/internal/api/http"
	"github.com/MohammadNE/PhoneBook/internal/repository"
	"github.com/MohammadNE/PhoneBook/pkg/logger"
	"github.com/MohammadNE/PhoneBook/pkg/rdbms"
	"github.com/MohammadNE/PhoneBook/pkg/token"
)

func Default() *Config {
	return &Config{
		Logger: &logger.Config{
			Development: true,
			Level:       "debug",
			Encoding:    "console",
		},
		HTTP: &http.Config{
			ListenPort: 8080,
		},
		RDBMS: &rdbms.Config{
			Host:     "localhost",
			Port:     5432,
			Username: "PHONEBOOK_USER",
			Password: "PHONEBOOK_PASSWORD",
			Database: "PHONEBOOK_DB",
		},
		Repository: &repository.Config{
			CursorSecret: "A?D(G-KaPdSgVkYp",
			Limit: struct {
				Min int "koanf:\"min\""
				Max int "koanf:\"max\""
			}{12, 48},
		},
		Token: &token.Config{
			PrivatePem: "-----BEGIN PRIVATE KEY-----\n" +
				"MC4CAQAwBQYDK2VwBCIEINyMNS8h9M9HO73Tg1BPr53p//qlqylO+wPKN8GrlsX7\n" +
				"-----END PRIVATE KEY-----",
			PublicPem: "-----BEGIN PUBLIC KEY-----\n" +
				"MCowBQYDK2VwAyEAqQsZ5iRNP3kdpNn3V/db9o/WkYHY8kkwQqCZGcDvJ+g=\n" +
				"-----END PUBLIC KEY-----",
			Expiration: 30 * time.Minute,
		},
	}
}
