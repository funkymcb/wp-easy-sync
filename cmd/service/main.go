package main

import (
	"cmd/service/main.go/pkg/config"
	"cmd/service/main.go/pkg/wordpress"
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

	// log.Println("Fetching Members from easyverein.com: ...")
	// easyvereinMembers, err := easyverein.GetMembers(client)
	// if err != nil {
	// 	log.Fatalln("Error fetching Members from easyverein.com, Error:", err)
	// }
	// log.Println("Fetching Members from easyverein.com: SUCCESS")

	// for i, member := range easyvereinMembers {
	// 	easyvereinMembers[i].LoginName = wordpress.GenerateLoginName(member)
	// }

	// log.Printf("Fetched %d Members from easyverein.com", len(easyvereinMembers))

	log.Println("Fetching Users from wordpress: ...")
	wordpressUsers, err := wordpress.GetUsers(client)
	if err != nil {
		log.Fatalf("Error fetching Users from %s, Error: %v",
			config.GetConfig().Wordpress.Host,
			err,
		)
	}
	log.Println("Fetching Users from wordpress: SUCCESS")

	log.Printf("Fetched %d Users from %s",
		len(wordpressUsers),
		config.GetConfig().Wordpress.Host,
	)
}

func printBanner() {
	fmt.Printf("\nWP-Easy-Sync:\n")
	fmt.Println("                         __o           o")
	fmt.Println("                       _ \\<_          <|/")
	fmt.Println("         ~~/\\O~^~~    (_)/(_)         / >")
	fmt.Printf("\nSyncs Members of Easyverein with Users of Wordpress\n")
	fmt.Printf("git: https://github.com/y-peter/wp-easy-sync\n\n")
}
