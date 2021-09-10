package main

import (
	"cmd/service/main.go/pkg/config"
	"cmd/service/main.go/pkg/easyverein"
	"flag"
	"fmt"
	"log"
)

func main() {
	configPathStr := flag.String("config", "./configs/config.yaml", "/path/to/config.yaml")
	flag.Parse()

	if err := config.LoadConfig(*configPathStr); err != nil {
		log.Fatalln("Config could not be loaded. Make sure to add the config.yaml to the specified path", err)
	}

	members, err := easyverein.ListMembers()
	if err != nil {
		log.Fatalln("Error fetching Members from easyverein, Error:", err)
	}

	for i, member := range members.Members {
		fmt.Printf("Member %d\n", i+1)
		fmt.Printf("\tFirst Name: %s\n", member.FirstName)
		fmt.Printf("\tLast Name: %s\n", member.FamilyName)
		fmt.Printf("\tEmail: %s\n", member.Email)
	}
}
