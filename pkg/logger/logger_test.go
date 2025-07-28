package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestInitLogger(t *testing.T) {
	tests := []struct {
		name        string
		env         string
		wantErr     bool
		expectPanic bool
		checkLogger func(t *testing.T, logger *zap.Logger)
	}{
		{
			name:        "local environment",
			env:         "local",
			wantErr:     false,
			expectPanic: false,
			checkLogger: func(t *testing.T, logger *zap.Logger) {
				assert.NotNil(t, logger, "Logger should not be nil")
				assert.Equal(t, zap.DebugLevel, logger.Level(), "Local env should use debug level")
			},
		},
		{
			name:        "development environment",
			env:         "development",
			wantErr:     false,
			expectPanic: false,
			checkLogger: func(t *testing.T, logger *zap.Logger) {
				assert.NotNil(t, logger, "Logger should not be nil")
				assert.Equal(t, zap.DebugLevel, logger.Level(), "Development env should use debug level")
			},
		},
		{
			name:        "staging environment",
			env:         "staging",
			wantErr:     false,
			expectPanic: false,
			checkLogger: func(t *testing.T, logger *zap.Logger) {
				assert.NotNil(t, logger, "Logger should not be nil")
				assert.Equal(t, zap.DebugLevel, logger.Level(), "Staging env should use debug level")
			},
		},
		{
			name:        "production environment",
			env:         "production",
			wantErr:     false,
			expectPanic: false,
			checkLogger: func(t *testing.T, logger *zap.Logger) {
				assert.NotNil(t, logger, "Logger should not be nil")
				assert.Equal(t, zap.InfoLevel, logger.Level(), "Production env should use info level")
			},
		},
		{
			name:        "unknown environment",
			env:         "unknown",
			wantErr:     true,
			expectPanic: true,
			checkLogger: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectPanic {
				assert.Panics(t, func() {
					_, _ = InitLogger(tt.env)
				}, "Expected panic for environment: %s", tt.env)
				return
			}

			logger, err := InitLogger(tt.env)
			if tt.wantErr {
				assert.Error(t, err, "Expected error for environment: %s", tt.env)
				return
			}

			assert.NoError(t, err, "Unexpected error for environment: %s", tt.env)
			tt.checkLogger(t, logger)

			assert.NotPanics(t, func() {
				logger.Info("test log message")
			}, "Logger should be able to log messages")

		})
	}
}
