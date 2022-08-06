package main

import (
	"fmt"

	"fizzbuzz-server/config"
	"fizzbuzz-server/http"

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

	panic(http.Start(conf))
}
