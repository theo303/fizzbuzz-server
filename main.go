package main

import (
	"fmt"

	"fizzbuzz-server/config"
	"fizzbuzz-server/http"
)

func main() {
	conf, errConf := config.InitEnvConf()
	if errConf != nil {
		panic(fmt.Errorf("error while initializing configuration: %w", errConf))
	}

	panic(http.Start(conf))
}
