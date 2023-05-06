package config

import (
	"github.com/MohammadNE/PhoneBook/pkg/logger"
	"github.com/MohammadNE/PhoneBook/pkg/rdbms"
)

// "github.com/CafeKetab/auth/pkg/crypto"
// "github.com/CafeKetab/auth/pkg/logger"
// "github.com/CafeKetab/auth/pkg/token"

type Config struct {
	Logger *logger.Config `koanf:"logger"`
	RDBMS  *rdbms.Config  `koanf:"rdbms"`
	// Token  *token.Config  `koanf:"token"`
	// Crypto *crypto.Config `koanf:"crypto"`
}
