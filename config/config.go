package config

import (
	"fmt"

	env "github.com/Netflix/go-env"
)

// Conf contains the program configuration
type Conf struct {
	Port int `env:"PORT,default=8080"`
}

// InitEnvConf initiate a Conf struct using env vars
func InitEnvConf() (Conf, error) {
	var conf Conf
	_, envErr := env.UnmarshalFromEnviron(&conf)
	if envErr != nil {
		return conf, fmt.Errorf("loading env vars: %w", envErr)
	}

	return conf, nil
}
