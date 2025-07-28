package cmd

import (
	"bytes"
	"context"
	"fmt"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestExecuteWithCommand(t *testing.T) {
	// Create a temporary context
	Execute()
	if rootCmd == nil {
		t.Fatal("Expected rootCmd to be initialized, got nil")
		ctx := context.Background()

		buffer := new(bytes.Buffer)
		rootCmd.SetOut(buffer)

		err := rootCmd.ExecuteContext(ctx)
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		assert.Contains(t, buffer.String(), "Test command executed", "Expected output to contain 'Test command executed'")
	}
}

func TestExecuteWithTestCommand(t *testing.T) {
	rootCmd := &cobra.Command{Use: "app"}

	testCmd := &cobra.Command{
		Use: "test",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintln(cmd.OutOrStdout(), "Test command executed")
		},
	}

	rootCmd.AddCommand(testCmd)

	buffer := new(bytes.Buffer)
	rootCmd.SetOut(buffer)
	rootCmd.SetErr(buffer)

	rootCmd.SetArgs([]string{"test"})

	err := rootCmd.Execute()

	assert.NoError(t, err)
	assert.Contains(t, buffer.String(), "Test command executed")
}
