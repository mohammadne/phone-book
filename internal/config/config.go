package config

import (
	"github.com/MohammadNE/PhoneBook/internal/api/http"
	"github.com/MohammadNE/PhoneBook/internal/repository"
	"github.com/MohammadNE/PhoneBook/pkg/logger"
	"github.com/MohammadNE/PhoneBook/pkg/rdbms"
	"github.com/MohammadNE/PhoneBook/pkg/token"
)

// "github.com/CafeKetab/auth/pkg/crypto"
// "github.com/CafeKetab/auth/pkg/logger"
// "github.com/CafeKetab/auth/pkg/token"

type Config struct {
	Logger     *logger.Config     `koanf:"logger"`
	HTTP       *http.Config       `koanf:"http"`
	RDBMS      *rdbms.Config      `koanf:"rdbms"`
	Repository *repository.Config `koanf:"repository"`
	Token      *token.Config      `koanf:"token"`
}
