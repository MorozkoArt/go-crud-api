package config

import (
    "time"
    "github.com/spf13/viper"
)

type Config struct {
    Server   ServerConfig   `mapstructure:"server"`
    Database DatabaseConfig `mapstructure:"database"`
    Auth     AuthConfig     `mapstructure:"auth"`
}

type ServerConfig struct {
    Port int `mapstructure:"port"`
}

type DatabaseConfig struct {
    Host     string `mapstructure:"host"`
    Port     int    `mapstructure:"port"`
    User     string `mapstructure:"user"`
    Password string `mapstructure:"password"`
    Name     string `mapstructure:"name"`
    SSLMode  string `mapstructure:"sslmode"`
}

type AuthConfig struct {
    JWTSecret   string        `mapstructure:"jwt_secret"`
    TokenExpiry time.Duration `mapstructure:"token_expiry"`
}

func LoadConfig() (*Config, error) {
    viper.SetConfigName("config")
    viper.SetConfigType("yaml")
    viper.AddConfigPath(".")
    
    viper.SetDefault("auth.token_expiry", "24h")
    
    if err := viper.ReadInConfig(); err != nil {
        return nil, err
    }

    var cfg Config
    if err := viper.Unmarshal(&cfg); err != nil {
        return nil, err
    }

    return &cfg, nil
}