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
	AddedUsers    map[int]string
	skipped       []Skip
	failed        []Skip
}

// TODO implement this
type Skip struct {
	SkipID   int
	Username string
	message  string
}

var syncStatus *Sync

var prefix = "[SYNC: %s] "

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
	prefixID := fmt.Sprintf(prefix, requestID)

	log.Println(prefixID, "################ RUNNING USER SYNC ################")
	client := resty.New()

	log.Println(prefixID, "Fetching Members from easyverein.com: ...")
	log.Println(prefixID, "(API can list max 100 users per page)")
	easyvereinMembers, err := easyverein.GetMembers(prefixID, client)
	if err != nil {
		msg := "Error fetching Members from easyverein.com, Error:"
		log.Println(prefixID, msg, err)
		syncStatus.Message = msg
	}
	log.Println(prefixID, "Fetching Members from easyverein.com: SUCCESS")

	for i, member := range *easyvereinMembers {
		(*easyvereinMembers)[i].LoginName = member.GenerateLoginName()
		(*easyvereinMembers)[i].Password = member.GeneratePassword(prefixID)
	}

	log.Printf("%s Fetched %d Members from easyverein.com", prefixID, len(*easyvereinMembers))

	log.Printf("%s Fetching Members from %s: ...", prefixID, config.GetConfig().Wordpress.Host)
	log.Println(prefixID, "(API can list max 100 users per page)")
	wordpressUsers, err := wordpress.GetUsers(prefixID, client)
	if err != nil {
		msg := fmt.Sprintf("%s error fetching users from %s, %v",
			prefixID,
			config.GetConfig().Wordpress.Host,
			err,
		)
		syncStatus.Message = msg
	}
	log.Printf("%s Fetching Users from %s: SUCCESS", prefixID, config.GetConfig().Wordpress.Host)

	log.Printf("%s Fetched %d Users from %s",
		prefixID,
		len(*wordpressUsers),
		config.GetConfig().Wordpress.Host,
	)

	log.Println(prefixID, "Running Synchronisation...")
	addedUsers := run(prefixID, client, easyvereinMembers, wordpressUsers)
	if err != nil {
		msg := fmt.Sprintf("%s Synchronisation of Users failed: %v", prefixID, err)
		syncStatus.Message = msg
	}

	addedUsersMap := make(map[int]string)
	for i, user := range addedUsers {
		name := fmt.Sprintf("%s %s", user.FirstName, user.LastName)
		addedUsersMap[i] = name
	}
	syncStatus.AddedUsers = addedUsersMap

	// clean things up
	client.SetCloseConnection(true)
	setSyncSuccess()
	cleanSyncArtifacts()
	// TODO HUGE BUG somehow setting slices to nil does not clear the length of the slice
	easyvereinMembers = nil
	wordpressUsers = nil
	log.Println(prefixID, "################ USER SYNC SUCCESSFUL ################")
	log.Println("keep listening for requests ...")
}

func newSync(ID string) {
	syncStatus.CurrentSyncID = ID
	syncStatus.Status = false
	syncStatus.Message = "sync running"
	syncStatus.AddedUsers = make(map[int]string)
}

// Run runs a synchronisation of two User slices.
// Every member that exists in easyverein but not in wordpress
// will be created in wp
// the users which were newly created will be returned
func run(prefix string, client *resty.Client, easyMembers, wpUsers *[]models.User) []models.User {
	// slice of Users to store the additional easyverein members
	var additionalUsers []models.User

EasyLoop:
	for _, easyMember := range *easyMembers {
		// check if loginname is on blacklist (config)
		for _, blacked := range config.GetConfig().Wordpress.Blacklist {
			if easyMember.LoginName == blacked {
				continue EasyLoop
			}
		}

		// skip empty easyverein members
		if easyMember.FirstName == "" {
			continue
		}

		// skip random 'Mustermann' members
		// don't know why easyverein returns these account but they are not shown in gui
		if strings.Contains(easyMember.FirstName, "Muster") ||
			strings.Contains(easyMember.LastName, "Muster") {
			continue
		}

		// loop over all wordpress users and check if loginname already exists
		for _, wpUser := range *wpUsers {
			// if loginnames are equal, the username already exists in wordpress
			// and we can continue with the next member of easyverein
			if compareLoginNames(easyMember.LoginName, wpUser.LoginName) {
				continue EasyLoop
			}
		}

		additionalUsers = append(additionalUsers, easyMember)
	}

	var counter, failed, skipped int
	for _, add := range additionalUsers {
		log.Printf("%s User %v has no Wordpress Account yet. Check if Account can be created",
			prefix,
			add,
		)

		if add.Email == "" {
			log.Printf("%s Skipping User %v: has no valid email Adress. Account can not be created.",
				prefix,
				add,
			)
			skipped++
			continue
		}

		err := wordpress.CreateUser(prefix, client, add)
		if err != nil {
			log.Println(err)
			failed++
			continue
		}
		counter++
	}
	if counter > 0 {
		log.Printf("%s Synchronisation done. %d WP-Accounts have been created. %d Accounts have been skipped. Creation of %d Accounts failed",
			prefix,
			counter,
			skipped,
			failed,
		)
	} else {
		log.Printf("%s Synchronisation done. No Accounts created. %d Accounts have been skipped. Creation of %d Accounts failed",
			prefix,
			skipped,
			failed,
		)
	}

	return additionalUsers
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

// since the loginnames in wordpress might be a big mess
// we have to check special cases
func compareLoginNames(n, m string) bool {
	// max.mustermann = max.mustermann
	if n == m {
		return true
	}

	// max.mustermann-zweitname = max.mustermann
	if strings.Split(n, "-")[0] == m {
		return true
	}

	//max.säge = max.saege
	//max.müller = max.mueller
	//max.möller = max.moeller
	//max.meißner = max.meissner
	//max-herber.mustermann = max.herbert.mustermann
	if strings.ReplaceAll(n, "ä", "ae") == m ||
		strings.ReplaceAll(n, "ü", "ue") == m ||
		strings.ReplaceAll(n, "ö", "oe") == m ||
		strings.ReplaceAll(n, "-", ".") == m ||
		strings.ReplaceAll(n, "ß", "ss") == m {
		return true
	}

	//max.müßnärü = max.muessnaerue
	allReplaced := n
	allReplaced = strings.Split(n, "-")[0]
	allReplaced = strings.ReplaceAll(allReplaced, "ä", "ae")
	allReplaced = strings.ReplaceAll(allReplaced, "ü", "ue")
	allReplaced = strings.ReplaceAll(allReplaced, "ö", "oe")
	allReplaced = strings.ReplaceAll(allReplaced, "ß", "ss")
	if allReplaced == m {
		return true
	}

	return false
}
