//go:build unit

package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadConfig(t *testing.T) {
	t.Run("should read config and return config instance", func(t *testing.T) {
		var err error

		err = os.Setenv(ServerPort, "8080")
		require.NoError(t, err)

		var cfg Config
		cfg, err = ReadConfig()
		defer os.Clearenv()

		assert.NoError(t, err)
		assert.IsType(t, &config{}, cfg)
	})

	t.Run("when read config server port parameter is not defined should return error", func(t *testing.T) {
		_, err := ReadConfig()

		assert.ErrorIs(t, err, EnvironmentVariablesNotDefined)
	})

}

func TestConfig_GetPactBrokerToken(t *testing.T) {
	cfg := &config{
		pactBrokerToken: "pact-broker-token",
	}
	token := cfg.GetPactBrokerToken()

	assert.Equal(t, token, cfg.pactBrokerToken)
}

func TestConfig_GetServerPort(t *testing.T) {
	cfg := &config{
		serverPort: "8080",
	}
	serverPort := cfg.GetServerPort()

	assert.Equal(t, serverPort, cfg.serverPort)
}
