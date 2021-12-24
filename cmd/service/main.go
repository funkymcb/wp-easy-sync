package main

import (
	"cmd/service/main.go/handler"
	"cmd/service/main.go/pkg/config"
	"flag"
	"fmt"
	"log"

	"github.com/savsgio/atreugo/v11"
)

func main() {
	configPathStr := flag.String("config", "./configs/config.yaml", "/path/to/config.yaml")
	flag.Parse()

	printBanner()

	if err := config.LoadConfig(*configPathStr); err != nil {
		log.Fatalln("Config could not be loaded. Make sure to add the config.yaml to the specified path", err)
	}

	apiServer := atreugo.New(config.GetConfig().API.AtreugoConfig())
	initAPIRoutes(apiServer)

	if err := apiServer.ListenAndServe(); err != nil {
		panic(err)
	}
}

func printBanner() {
	fmt.Printf("\nWP-Easy-Sync-API:\n")
	fmt.Println("                         __o           o")
	fmt.Println("                       _ \\<_          <|/")
	fmt.Println("         ~~/\\O~^~~    (_)/(_)         / >")
	fmt.Printf("\nAPI to handle user synchronisation between easyverein and wordpress\n")
	fmt.Printf("git: https://github.com/y-peter/wp-easy-sync\n\n")
}

func initAPIRoutes(server *atreugo.Atreugo) {
	server.GET("/sync", handler.SyncEasyToWP)
	server.GET("/sync/status/{requestID}", handler.SyncStatus)
	server.GET("/sync/csv", handler.SyncCSVToEasy)
}
