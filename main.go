package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"fizzbuzz-server/api"
	"fizzbuzz-server/config"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
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

	api := api.Init(conf)

	go func() {
		if errServ := api.Run(); errServ != nil {
			log.Warn().Err(errServ).Msg("server exited")
		}
	}()

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt)

	<-stopChan

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	api.Shutdown(ctx)
}
