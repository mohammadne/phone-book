package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/MohammadNE/PhoneBook/cmd"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func main() {
	const description = "PhoneBook application"
	root := &cobra.Command{Short: description}

	trap := make(chan os.Signal, 1)
	signal.Notify(trap, syscall.SIGINT, syscall.SIGTERM)

	root.AddCommand(
		cmd.Server{}.Command(trap),
		cmd.Migrate{}.Command(trap),
	)

	if err := root.Execute(); err != nil {
		log.Fatal("failed to execute root command", zap.Error(err))
	}
}
