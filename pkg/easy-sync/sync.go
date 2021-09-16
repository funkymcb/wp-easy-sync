package easysync

import (
	"cmd/service/main.go/pkg/models"
	"fmt"
	"strings"
)

// Run() runs a synchronisation of two User slices.
// Every member that exists in easyverein but not in wordpress
// will be created in wp
func Run(easyMembers, wpUsers []models.User) error {
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
			if easyMember.LoginName == wpUser.LoginName {
				continue EasyLoop
			}
		}
		additionalUsers = append(additionalUsers, easyMember)
	}

	for i, add := range additionalUsers {
		fmt.Printf("%d. User will be created\n\tName: %s %s\n\tUsername: %s\n",
			i,
			add.FirstName,
			add.LastName,
			add.LoginName,
		)
	}

	return nil
}
