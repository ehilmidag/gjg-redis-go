//go:build unit

package config

//import (
//	"os"
//	"testing"
//
//	"github.com/stretchr/testify/assert"
//	"github.com/stretchr/testify/require"
//)
//
//func TestReadConfig(t *testing.T) {
//	t.Run("should read config and return config instance", func(t *testing.T) {
//		var err error
//
//		err = os.Setenv(ServerPort, "8080")
//		require.NoError(t, err)
//		mySQLConf, err := ReadMySqlConfig()
//		redisConf, err := ReadRedisConfig()
//		var cfg Config
//		cfg, err = ReadConfig()
//		defer os.Clearenv()
//
//		assert.NoError(t, err)
//		assert.IsType(t, &Config{
//			ServerPort: "8080",
//			mySQL:      redisConf,
//			Redis:      mySQLConf,
//		}, cfg)
//	})
//
//	t.Run("when read config server port parameter is not defined should return error", func(t *testing.T) {
//		_, err := ReadConfig()
//
//		assert.ErrorIs(t, err, EnvironmentVariablesNotDefined)
//	})
//
//}
//func TestConfig_ReadRedisConfig(t *testing.T) {
//	Rediscfg := &RedisConfig{
//		Addr:     "localhost:6379",
//		Password: "",
//		Db:       0,
//	}
//	mysqlConfig := &MySQLConfig{
//		Username: "root",
//		Password: "",
//		Hostname: "127.0.0.1:3306",
//		Dbname:   "user",
//	}
//}
//
//func TestConfig_GetServerPort(t *testing.T) {
//	Rediscfg := &RedisConfig{
//		Addr:     "localhost:6379",
//		Password: "",
//		Db:       0,
//	}
//	mysqlConfig := &MySQLConfig{
//		Username: "root",
//		Password: "",
//		Hostname: "127.0.0.1:3306",
//		Dbname:   "user",
//	}
//	cfg := &Config{
//		serverPort: "8080",
//		MySQL:      mysqlConfig,
//		Redis:      Rediscfg,
//	}
//	a, err := ReadConfig()
//	a.MySQL
//
//	assert.Equal(t, serverPort, cfg.serverPort)
//}
