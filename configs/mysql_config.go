package configs

type MySQLConfig struct {
	DSN string `json:"dsn" yaml:"dsn" mapstructure:"dsn"`
}
