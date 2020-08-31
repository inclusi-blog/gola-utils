package model

type DBConnectionPoolConfig struct {
	MaxOpenConnections             int `json:"maxOpenConnections"`
	MaxIdleConnections             int `json:"maxIdleConnections"`
	MaxConnectionLifetimeInMinutes int `json:"maxConnectionLifetimeInMinutes"`
}
