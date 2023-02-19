package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

//PubKeyPem can contain Pem-encoded public key
type Config struct {
	Debug               bool     `env:"DEBUG" env-default:"false"`
	Host                string   `env:"HOST" env-default:"0.0.0.0"`
	Port                int      `env:"PORT" env-default:"8080"`
	Prefix              string   `env:"PREFIX" env-default:"rigel"`
	Version             string   `env:"VERSION" env-default:"0.6.0"`
	CORS                bool     `env:"CORS" env-default:"false"`
	AllowMethods        []string `env:"AllowMethods" env-default:"GET,POST"`
	AllowOrigins        []string `env:"AllowOrigins" env-default:"*"`
	AllowHeaders        []string `env:"AllowHeaders" env-default:"Accept,Content-Type"`
	RedisAddress        string   `env:"REDIS_ADDRESS" env-default:"localhost:6379"`
	RedisPassword       string   `env:"REDIS_PASSWORD" env-default:""`
	RedisDB             int      `env:"REDIS_DB" env-default:"0"`
	RedisTimeout        int      `env:"REDIS_TIMEOUT" env-default:"3"`
	RedisExpiration     int      `env:"REDIS_EXPIRATION" env-default:"2592000"` // Expiration date for redis | default = 30 * 24 * 60 * 60
	Cap                 int      `env:"CAP" env-default:"500"`                  // Capacity for lfu cache
	SignatureValidation bool     `env:"SIGNATURE_VALIDATION" env-default:"true"`
	XKey                string   `env:"X_KEY" env-default:"secretkey"`
	XSalt               string   `env:"X_SALT" env-default:"secretsalt"`
}

func GetConfig() *Config {
	var cfg Config
	cleanenv.ReadEnv(&cfg)
	return &cfg
}
