package configs

type MySQLConfig struct {
	DSN string     `json:"dsn" yaml:"dsn"`
	Log *LogConfig `json:"log" yaml:"log"`
}

func GetMySQLConfig() *MySQLConfig {
	return c.MySQL
}
