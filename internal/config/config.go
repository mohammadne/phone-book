package config

import (
	"github.com/MohammadNE/PhoneBook/internal/repository"
	"github.com/MohammadNE/PhoneBook/pkg/logger"
	"github.com/MohammadNE/PhoneBook/pkg/rdbms"
	"github.com/MohammadNE/PhoneBook/pkg/token"
)

type Config struct {
	Logger     *logger.Config     `koanf:"logger"`
	RDBMS      *rdbms.Config      `koanf:"rdbms"`
	Repository *repository.Config `koanf:"repository"`
	Token      *token.Config      `koanf:"token"`
}
