package main

import (
	"cmd/service/main.go/pkg/config"
	"cmd/service/main.go/pkg/easyverein"
	"flag"
	"fmt"
	"log"

	"github.com/go-resty/resty/v2"
)

func main() {
	configPathStr := flag.String("config", "./configs/config.yaml", "/path/to/config.yaml")
	flag.Parse()

	printBanner()

	if err := config.LoadConfig(*configPathStr); err != nil {
		log.Fatalln("Config could not be loaded. Make sure to add the config.yaml to the specified path", err)
	}

	client := resty.New()

	log.Println("Fetching Members from easyverein: ...")
	easyvereinMembers, err := easyverein.GetMembers(client)
	if err != nil {
		log.Fatalln("Error fetching Members from easyverein, Error:", err)
	}
	log.Println("Fetching Members from easyverein: SUCCESS")

	log.Printf("Fetched %d Members from Easyverein", len(easyvereinMembers))
}

func printBanner() {
	fmt.Printf("\nWVC-Sync:\n")
	fmt.Println("                         __o           o")
	fmt.Println("                       _ \\<_          <|/")
	fmt.Println("         ~~/\\O~^~~    (_)/(_)         / >")
	fmt.Printf("\nSyncs easyverein Members with Members of Wordpress\n")
	fmt.Printf("git: https://github.com/y-peter/wvc-sync\n\n")
}
