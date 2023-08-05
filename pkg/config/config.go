package config

import (
	"os"
)

type Config interface {
	GetServerPort() string
	GetPactBrokerURl() string
	GetPactBrokerToken() string
}

type config struct {
	serverPort      string
	pactBrokerURl   string
	pactBrokerToken string
}

func ReadConfig() (Config, error) {
	serverPort := os.Getenv(ServerPort)
	if serverPort == "" {
		return nil, EnvironmentVariablesNotDefined
	}

	pactBrokerURl := os.Getenv(PactBrokerURL)
	if serverPort == "" {
		return nil, EnvironmentVariablesNotDefined
	}

	pactBrokerToken := os.Getenv(PactBrokerToken)
	if pactBrokerToken == "" {
		return nil, EnvironmentVariablesNotDefined
	}

	return &config{
		serverPort:      serverPort,
		pactBrokerURl:   pactBrokerURl,
		pactBrokerToken: pactBrokerToken,
	}, nil
}

func (config *config) GetServerPort() string {
	return config.serverPort
}

func (config *config) GetPactBrokerURl() string {
	return config.pactBrokerURl
}

func (config *config) GetPactBrokerToken() string {
	return config.pactBrokerToken
}
