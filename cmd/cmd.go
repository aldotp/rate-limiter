package cmd

import (
	"context"
	"fmt"

	"github.com/aldotp/rate-limiter/cmd/http"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{}
)

func init() {
	ctx := context.Background()

	restCmd := cobra.Command{
		Use:   "rest",
		Short: "Rest is a command to start Restful Api server",
		Run: func(cmd *cobra.Command, args []string) {
			http.RunHTTPServer(ctx)
		},
	}

	rootCmd.AddCommand(&restCmd)
}

func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		return fmt.Errorf("error executing root command: %w", err)
	} else {
		return nil
	}
}
