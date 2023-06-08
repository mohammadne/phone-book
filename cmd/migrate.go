package cmd

import (
	"os"

	"github.com/mohammadne/phone-book/internal/config"
	"github.com/mohammadne/phone-book/internal/models"
	"github.com/mohammadne/phone-book/internal/repository"
	"github.com/mohammadne/phone-book/pkg/logger"
	"github.com/mohammadne/phone-book/pkg/rdbms"
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
		logger.Fatal("Invalid arguments given", zap.Any("args", args))
	}

	rdbms, err := rdbms.New(cfg.RDBMS)
	if err != nil {
		logger.Fatal("Error creating rdbms", zap.Error(err))
	}

	repository := repository.New(logger, cfg.Repository, rdbms)
	if err := repository.Migrate(models.Migrate(args[0])); err != nil {
		logger.Fatal("Error migrating", zap.String("migration", args[0]), zap.Error(err))
	}

	logger.Info("Database has been migrated successfully", zap.String("migration", args[0]))
}
