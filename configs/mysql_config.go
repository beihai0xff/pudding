package configs

type MySQLConfig struct {
	DSN string     `json:"dsn" yaml:"dsn" mapstructure:"dsn"`
	Log *LogConfig `json:"log" yaml:"log" mapstructure:"log"`
}

func GetMySQLConfig() *MySQLConfig {
	return c.MySQL
}
