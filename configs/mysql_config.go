package configs

type MySQLConfig struct {
	DSN string `json:"dsn" yaml:"dsn"`
}

func GetMySQLConfig() *MySQLConfig {
	return c.MySQL
}
