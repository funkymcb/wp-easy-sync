package main

import (
	"cmd/service/main.go/pkg/config"
	"flag"
	"log"
)

func main() {
	configPathStr := flag.String("config", "./configs/config.yaml", "/path/to/config.yaml")
	flag.Parse()

	if err := config.LoadConfig(*configPathStr); err != nil {
		log.Fatalln("Config could not be loaded. Make sure to add the config.yaml to the specified path", err)
	}
}
