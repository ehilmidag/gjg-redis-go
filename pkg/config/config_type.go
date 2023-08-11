package config

var (
	EnvironmentVariableNotDefined = "%s variable is not defined"
)

// #nosec
const (
	IsAtRemote = "IS_AT_REMOTE"

	ServerPort    = "SERVER_PORT"
	MySQLUsername = "MYSQL_USER"
	MySQLPassword = "MYSQL_PASSWORD"
	MySQLHostname = "MYSQL_HOST"
	MySQLDbName   = "MYSQL_DB_NAME"
	RedisAddr     = "REDIS_ADDR"
	RedisPassword = "REDIS_PASSWORD"
	RedisDb       = "REDIS_DB"
)

type MySQLConfig struct {
	Username string
	Password string
	Hostname string
	Dbname   string
}

type RedisConfig struct {
	Addr     string
	Password string
	Db       int
}
