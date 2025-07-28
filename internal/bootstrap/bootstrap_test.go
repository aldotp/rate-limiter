package bootstrap_test

import (
	"context"
	"testing"

	"github.com/aldotp/rate-limiter/internal/bootstrap"
	"github.com/stretchr/testify/assert"
)

func TestNewBootstrap(t *testing.T) {
	ctx := context.Background()
	b := bootstrap.NewBootstrap(ctx)

	assert.NotNil(t, b)
	assert.Nil(t, b.RedisClient)
	assert.Nil(t, b.Log)
	assert.Nil(t, b.Config)
}

func TestBuildDependencies(t *testing.T) {
	tests := []struct {
		name        string
		setup       func(t *testing.T) *bootstrap.Bootstrap
		expectError bool
	}{
		{
			name: "successful dependency build",
			setup: func(t *testing.T) *bootstrap.Bootstrap {
				b := bootstrap.NewBootstrap(context.Background())
				return b
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := tt.setup(t)
			result, err := b.BuildDependencies()

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.NotNil(t, result.Config)
				assert.NotNil(t, result.Log)
				assert.NotNil(t, result.RedisClient)
			}
		})
	}
}
