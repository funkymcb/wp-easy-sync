package handler

import (
	"cmd/service/main.go/pkg/config"
	"cmd/service/main.go/pkg/easyverein"
	"cmd/service/main.go/pkg/wordpress"
	"fmt"
	"log"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/savsgio/atreugo/v11"
)

// SyncEasyToWP triggers the sync between Easyverein and Wordpress
// Method: GET
// Responses: 200, TODO
func SyncEasyToWP(ctx *atreugo.RequestCtx) error {
	client := resty.New()

	log.Println("Fetching Members from easyverein.com: ...")
	log.Println("(API can list max 100 users per page)")
	easyvereinMembers, err := easyverein.GetMembers(client)
	if err != nil {
		log.Fatalln("Error fetching Members from easyverein.com, Error:", err)
		ctx.Error(err.Error(), http.StatusInternalServerError)
		return fmt.Errorf(err.Error())
	}
	log.Println("Fetching Members from easyverein.com: SUCCESS")

	for i, member := range easyvereinMembers {
		easyvereinMembers[i].LoginName = member.GenerateLoginName()
		easyvereinMembers[i].Password = member.GeneratePassword()
	}

	log.Printf("Fetched %d Members from easyverein.com", len(easyvereinMembers))

	log.Printf("Fetching Members from %s: ...", config.GetConfig().Wordpress.Host)
	log.Println("(API can list max 100 users per page)")
	wordpressUsers, err := wordpress.GetUsers(client)
	if err != nil {
		log.Printf("error fetching users from %s, %v",
			config.GetConfig().Wordpress.Host,
			err,
		)
		ctx.Error(err.Error(), http.StatusInternalServerError)
		return fmt.Errorf(err.Error())
	}
	log.Printf("Fetching Users from %s: SUCCESS", config.GetConfig().Wordpress.Host)

	log.Printf("Fetched %d Users from %s",
		len(wordpressUsers),
		config.GetConfig().Wordpress.Host,
	)

	// TODO optimize user handling first
	// log.Println("Running Synchronisation...")
	// if err = easysync.Run(client, easyvereinMembers, wordpressUsers); err != nil {
	// 	log.Fatalln("Synchronisation of Users failed:", err)
	// }
	return ctx.JSONResponse("sync successfull")
}
