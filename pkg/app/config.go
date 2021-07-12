package app

import (
	"github.com/spf13/viper"
)

const (
	DefaultGRPCHost      = "0.0.0.0"
	DefaultGRPCPort      = 8080
	DefaultGatewayHost   = "0.0.0.0"
	DefaultGatewayPort   = 8090
	DefaultStoreHost     = "localhost"
	DefaultStorePort     = 5432
	DefaultStoreUser     = "postgres"
	DefaultStorePassword = "postgres"
	DefaultStoreDatabase = "postgres"
)

type Config struct {
	GRPCHost      string
	GRPCPort      int
	GatewayHost   string
	GatewayPort   int
	StoreHost     string
	StorePort     int
	StoreUser     string
	StorePassword string
	StoreDatabase string
}

func LoadConfig() *Config {
	v := viper.New()
	v.SetDefault("grpc_host", DefaultGRPCHost)
	v.SetDefault("grpc_port", DefaultGRPCPort)
	v.SetDefault("gateway_host", DefaultGatewayHost)
	v.SetDefault("gateway_port", DefaultGatewayPort)
	v.SetDefault("store_host", DefaultStoreHost)
	v.SetDefault("store_port", DefaultStorePort)
	v.SetDefault("store_user", DefaultStoreUser)
	v.SetDefault("store_password", DefaultStorePassword)
	v.SetDefault("store_database", DefaultStoreDatabase)

	v.BindEnv("grpc_host")
	v.BindEnv("grpc_port")
	v.BindEnv("gateway_host")
	v.BindEnv("gateway_port")
	v.BindEnv("store_host")
	v.BindEnv("store_port")
	v.BindEnv("store_user")
	v.BindEnv("store_password")
	v.BindEnv("store_database")

	config := Config{
		GRPCHost:      v.GetString("grpc_host"),
		GRPCPort:      v.GetInt("grpc_port"),
		GatewayHost:   v.GetString("gateway_host"),
		GatewayPort:   v.GetInt("gateway_port"),
		StoreHost:     v.GetString("store_host"),
		StorePort:     v.GetInt("store_port"),
		StoreUser:     v.GetString("store_user"),
		StorePassword: v.GetString("store_password"),
		StoreDatabase: v.GetString("store_database"),
	}
	return &config
}
