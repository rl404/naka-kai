package main

import (
	"github.com/rl404/naka-kai/internal/utils"
	"github.com/spf13/cobra"
)

func main() {
	cmd := cobra.Command{
		Use:   "naka-kai",
		Short: "Naka-kai Discord Bot",
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "bot",
		Short: "Run bot",
		RunE: func(*cobra.Command, []string) error {
			return bot()
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "migrate",
		Short: "Run migration",
		RunE: func(*cobra.Command, []string) error {
			return migrate()
		},
	})

	if err := cmd.Execute(); err != nil {
		utils.Fatal(err.Error())
	}
}
