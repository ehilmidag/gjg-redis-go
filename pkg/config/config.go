package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	ServerPort string
	Redis      *RedisConfig
}

func ReadConfig() (*Config, error) {
	serverPort := os.Getenv(ServerPort)
	if serverPort == "" {
		serverPort = "8080"
	}

	redisConfig, err := ReadRedisConfig()
	if err != nil {
		return nil, err
	}
	return &Config{
		ServerPort: serverPort,
		Redis:      redisConfig,
	}, nil
}

func ReadRedisConfig() (*RedisConfig, error) {
	redisAddr := os.Getenv(RedisAddr)
	if redisAddr == "" {
		return nil, fmt.Errorf(EnvironmentVariableNotDefined, redisAddr)

	}
	redisPassword := os.Getenv(RedisPassword)

	redisDB := os.Getenv(RedisDb)
	redisDBint, err := strconv.Atoi(redisDB)
	if err != nil {
		fmt.Println("Hata:", err)
		return nil, err
	}
	return &RedisConfig{
		Addr:     redisAddr,
		Password: redisPassword,
		Db:       redisDBint,
	}, nil
}
