package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	ServerPort string
	mySQL      *MySQLConfig
	Redis      *RedisConfig
}

func ReadConfig() (*Config, error) {
	serverPort := os.Getenv(ServerPort)
	if serverPort == "" {
		serverPort = "8080"
	}

	mySQLconfig, err := ReadMySqlConfig()
	if err != nil {
		return nil, err
	}
	redisConfig, err := ReadRedisConfig()
	if err != nil {
		return nil, err
	}
	return &Config{
		ServerPort: serverPort,
		mySQL:      mySQLconfig,
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

func ReadMySqlConfig() (*MySQLConfig, error) {

	mySqlUsername := os.Getenv(MySQLUsername)
	if mySqlUsername == "" {
		return nil, fmt.Errorf(EnvironmentVariableNotDefined, MySQLUsername)
	}
	mySqlPassword := os.Getenv(MySQLPassword)

	mySQLHostname := os.Getenv(MySQLHostname)
	if mySQLHostname == "" {
		return nil, fmt.Errorf(EnvironmentVariableNotDefined, MySQLHostname)

	}
	mySQLDbName := os.Getenv(MySQLDbName)
	if mySQLDbName == "" {
		return nil, fmt.Errorf(EnvironmentVariableNotDefined, MySQLDbName)
	}

	return &MySQLConfig{
		Username: mySqlUsername,
		Password: mySqlPassword,
		Hostname: mySQLHostname,
		Dbname:   mySQLDbName,
	}, nil

}
