package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

//PubKeyPem can contain Pem-encoded public key
type Config struct {
	Debug         bool   `env:"DEBUG" env-default:"false"`
	Host          string `env:"HOST" env-default:"0.0.0.0"`
	Port          int    `env:"PORT" env-default:"8080"`
	Prefix        string `env:"PREFIX" env-default:"rigel"`
	Version       string `env:"VERSION" env-default:"1.0.0"`
	RedisAddress  string `env:"REDIS_ADDRESS" env-default:"localhost:6379"`
	RedisPassword string `env:"REDIS_PASSWORD" env-default:""`
	RedisDB       int    `env:"REDIS_DB" env-default:"0"`
	Cap           int    `env:"CAP" env-default:"500"`
	Alg           string `env:"ALG" env-default:"RS256"`
	PubKeyPem     []byte `env:"PUBLIC_KEY_PEM"`
}

func GetConfig() *Config {
	var cfg Config
	cleanenv.ReadEnv(&cfg)
	return &cfg
}
