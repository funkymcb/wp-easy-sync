package easysync

import (
	"cmd/service/main.go/pkg/config"
	"cmd/service/main.go/pkg/easyverein"
	"cmd/service/main.go/pkg/models"
	"cmd/service/main.go/pkg/wordpress"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

type Sync struct {
	CurrentSyncID string
	LastSync      time.Time
	Status        bool
	Message       string
	CreatedUsers  map[int]string
}

var syncStatus *Sync

func init() {
	syncStatus = &Sync{
		Status:  false,
		Message: "no sync has been triggered yet",
	}
}

// GetSyncStatus returns the sync object
// if the requestIDs of the current (or previous) syncs are matching
func GetSyncStatus(requestID string) (*Sync, error) {
	if requestID == syncStatus.CurrentSyncID {
		return syncStatus, nil
	}
	return nil, fmt.Errorf("no active sync with this request id found")
}

// TriggerSync triggers a user synchronisation
// a new sync object will be initialized with success set to false
// if sync was successful the value will be set to true
func TriggerSync(requestID string) {
	newSync(requestID)

	log.Println("################ RUNNING USER SYNC ################")
	client := resty.New()

	log.Println("Fetching Members from easyverein.com: ...")
	log.Println("(API can list max 100 users per page)")
	easyvereinMembers, err := easyverein.GetMembers(client)
	if err != nil {
		msg := "Error fetching Members from easyverein.com, Error:"
		log.Println(msg, err)
		syncStatus.Message = msg
	}
	log.Println("Fetching Members from easyverein.com: SUCCESS")

	for i, member := range *easyvereinMembers {
		(*easyvereinMembers)[i].LoginName = member.GenerateLoginName()
		(*easyvereinMembers)[i].Password = member.GeneratePassword()
	}

	log.Printf("Fetched %d Members from easyverein.com", len(*easyvereinMembers))

	log.Printf("Fetching Members from %s: ...", config.GetConfig().Wordpress.Host)
	log.Println("(API can list max 100 users per page)")
	wordpressUsers, err := wordpress.GetUsers(client)
	if err != nil {
		msg := fmt.Sprintf("error fetching users from %s, %v",
			config.GetConfig().Wordpress.Host,
			err,
		)
		syncStatus.Message = msg
	}
	log.Printf("Fetching Users from %s: SUCCESS", config.GetConfig().Wordpress.Host)

	log.Printf("Fetched %d Users from %s",
		len(*wordpressUsers),
		config.GetConfig().Wordpress.Host,
	)

	// TODO optimize user handling first
	// log.Println("Running Synchronisation...")
	// if err = run(client, easyvereinMembers, wordpressUsers); err != nil {
	//  msg := fmt.Sprintf("Synchronisation of Users failed: %v", err)
	//  syncStatus.Message= msg
	// }

	// clean things up
	client.SetCloseConnection(true)
	setSyncSuccess()
	cleanSyncArtifacts()
	easyvereinMembers = nil
	wordpressUsers = nil
	log.Println("################ USER SYNC SUCCESSFUL ################")
	log.Println("keep listening for requests ...")
}

func newSync(ID string) {
	syncStatus.CurrentSyncID = ID
	syncStatus.Status = false
	syncStatus.Message = "sync running"
}

// Run runs a synchronisation of two User slices.
// Every member that exists in easyverein but not in wordpress
// will be created in wp
// the users which were newly created will be returned
func run(client *resty.Client, easyMembers, wpUsers []models.User) ([]models.User, error) {
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
			return nil, fmt.Errorf("user could not be created: %v", err)
		}
		counter = i
	}
	if counter > 0 {
		log.Printf("Synchronisation done. %d WP-Accounts have been created", counter)
	} else {
		log.Println("Synchronisation done. Everything up to date. No Accounts created")
	}

	return additionalUsers, nil
}

func setSyncSuccess() {
	syncStatus.Status = true
	syncStatus.LastSync = time.Now()
	syncStatus.Message = "Sync successful"
}

func cleanSyncArtifacts() {
	easyverein.Page = 1
	wordpress.Page = 1
}
