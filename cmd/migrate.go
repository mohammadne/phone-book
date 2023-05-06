package cmd

import (
	"os"

	"github.com/MohammadNE/PhoneBook/internal/config"
	"github.com/MohammadNE/PhoneBook/internal/models"
	"github.com/MohammadNE/PhoneBook/internal/repository"
	"github.com/MohammadNE/PhoneBook/pkg/logger"
	"github.com/MohammadNE/PhoneBook/pkg/rdbms"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type Migrate struct{}

func (m Migrate) Command(trap chan os.Signal) *cobra.Command {
	run := func(_ *cobra.Command, args []string) {
		m.main(config.Load(true), args, trap)
	}

	return &cobra.Command{
		Use:       "migrate",
		Short:     "run migrations",
		Run:       run,
		Args:      cobra.OnlyValidArgs,
		ValidArgs: []string{"up", "down"},
	}
}

func (m *Migrate) main(cfg *config.Config, args []string, trap chan os.Signal) {
	logger := logger.NewZap(cfg.Logger)

	if len(args) != 1 {
		logger.Fatal("invalid arguments given", zap.Any("args", args))
	}

	rdbms, err := rdbms.New(cfg.RDBMS)
	if err != nil {
		logger.Fatal("Error creating rdbms", zap.Error(err))
	}

	repository := repository.New(logger, rdbms)
	if args[0] == "up" {
		if err := repository.Migrate(models.Up); err != nil {
			logger.Fatal("Error migrating up", zap.Error(err))
		}
	} else {
		if err := repository.Migrate(models.Down); err != nil {
			logger.Fatal("Error migrating down", zap.Error(err))
		}
	}

	logger.Info("Database has been migrated successfully", zap.String("migration", args[0]))
}
