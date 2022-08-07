package main

import (
	"fmt"

	"fizzbuzz-server/api"
	"fizzbuzz-server/config"

	"github.com/rs/zerolog"
)

func main() {
	conf, errConf := config.InitEnvConf()
	if errConf != nil {
		panic(fmt.Errorf("error while initializing configuration: %w", errConf))
	}

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	logLevel, errLvl := zerolog.ParseLevel(conf.LogLevel)
	if errLvl != nil {
		panic(fmt.Errorf("error while parsing log level: %w", errLvl))
	}
	zerolog.SetGlobalLevel(logLevel)

	panic(api.Start(conf))
}
