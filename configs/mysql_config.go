// Package configs provides config management
package configs

// MySQLConfig MySQLConfig Config
type MySQLConfig struct {
	// DSN is the data source name
	DSN string `json:"dsn" yaml:"dsn" mapstructure:"dsn"`
}
