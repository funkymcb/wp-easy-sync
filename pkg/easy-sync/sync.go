package easysync

import (
	"cmd/service/main.go/pkg/models"
	"cmd/service/main.go/pkg/wordpress"
	"fmt"
	"log"
	"strings"

	"github.com/go-resty/resty/v2"
)

// Run() runs a synchronisation of two User slices.
// Every member that exists in easyverein but not in wordpress
// will be created in wp
func Run(client *resty.Client, easyMembers, wpUsers []models.User) error {
	// slice of Users to store the additional easyverein members
	var additionalUsers []models.User

EasyLoop:
	for _, easyMember := range easyMembers {
		// skip empty easyverein members
		if easyMember.FirstName == "" {
			continue
		}

		// skip random 'Mustermann' members
		if strings.Contains(easyMember.FirstName, "Muster") ||
			strings.Contains(easyMember.LastName, "Muster") {
			continue
		}

		for _, wpUser := range wpUsers {
			fullWPUserName := fmt.Sprintf("%s %s",
				wpUser.FirstName,
				wpUser.LastName,
			)
			if easyMember.Email == wpUser.Email {
				continue EasyLoop
			}
			if easyMember.LoginName == wpUser.LoginName {
				continue EasyLoop
			}
			if strings.Contains(fullWPUserName, easyMember.FirstName) &&
				strings.Contains(fullWPUserName, easyMember.LastName) {
				continue EasyLoop
			}
		}
		additionalUsers = append(additionalUsers, easyMember)
	}

	var counter int
	for i, add := range additionalUsers {
		log.Printf("User %s %s has no Wordpress Account yet.",
			add.FirstName,
			add.LastName,
		)

		if add.Email == "" {
			log.Println("User has no valid email Adress. Account can not be created. Skipping")
			continue
		}

		log.Printf("Account %s will be created", add.LoginName)
		err := wordpress.CreateUser(client, add)
		if err != nil {
			return fmt.Errorf("User could not be created: %v", err)
		}
		counter = i
	}
	if counter > 0 {
		log.Printf("Synchronisation done. %d WP-Accounts have been created", counter)
	} else {
		log.Println("Synchronisation done. Everything up to date. No Accounts created")
	}

	return nil
}
