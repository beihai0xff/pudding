package configs

// MySQLConfig MySQL Config
type MySQLConfig struct {
	DSN string `json:"dsn" yaml:"dsn" mapstructure:"dsn"`
}
