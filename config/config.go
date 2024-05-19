package config

import (
	"time"
)

const (
	RedisTTL = 60 * time.Second

	RedisConfigAddr     = "localhost:6379"
	RedisConfigPassword = ""
	RedisConfigDB       = 0

	PostgresConfigAddr     = ":5432"
	PostgresConfigUser     = "postgres"
	PostgresConfigPassword = "qwertfdsa"
	PostgreConfigDB        = "users"

	HTTPServerPort = 8080
)
