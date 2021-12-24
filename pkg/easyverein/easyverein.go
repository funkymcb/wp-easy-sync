package easyverein

import (
	"cmd/service/main.go/pkg/config"
	"cmd/service/main.go/pkg/models"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/go-resty/resty/v2"
)

// EasyVereinResponse stores all necessarry values of API response
type EasyVereinResponse struct {
	// the url for the next page, null if last page
	Next    string        `json:"next"`
	Members []models.User `json:"results"`
}

// EasyVereinMember represents the struct of a member as exported from
// the local csv export
type EasyVereinMember struct {
	MemberID          string
	Salution          string
	FirstName         string
	LastName          string
	AdditionalAddress string
	Street            string
	PostalCode        string
	Location          string
	Phone             string
	Country           string
	IBAN              string
	BIC               string
	Birthdate         string
	PhoneBiz          string
	Fax               string
	Mobile            string
	EMail             string
	Nationality       string
	Gender            string
	MaritalState      string
	Job               string
	Status            string
	Bank              string
	EntryDate         string
	LeavingDate       string
	Department        string
	Roles             string
	MandateID         string
	Title             string
}

var prefix = "[SYNC: %s] "

var Page = 1
var members *[]models.User

func init() {
	members = &[]models.User{}
}

// GetMembers() unmarshals the API response of the contact-details endpoint
// into a slice of Users
func GetMembers(prefix string, client *resty.Client) (*[]models.User, error) {
	log.Printf("%s Fetching easyverein members from page: %d", prefix, Page)
	var easyResponse EasyVereinResponse

	// requestURI = https://easyverein.com/api/stable/contact-details?limit100&page=%d
	requestURI := config.GetConfig().Easyverein.APIRequestURI("contact-details", Page)

	resp, err := makeAPIRequest(client, requestURI)
	if err != nil {
		return members, fmt.Errorf("%s could not perform get request to easyverein contact-details endpoint: %v", prefix, err)
	}

	err = json.Unmarshal(resp.Body(), &easyResponse)
	if err != nil {
		return members, err
	}

	// call GetMembers() recursively until no next page
	if easyResponse.Next != "" {
		*members = append(*members, easyResponse.Members...)
		Page += 1
		GetMembers(prefix, client)
	} else {
		// append members of the last page
		*members = append(*members, easyResponse.Members...)
	}

	return members, nil
}

// SyncCSV takes the data from POST request (base64 encoded csv file)
// and syncs the members listed with easyverein
// makes API POST request using resty
// headers:
//   "Authorization": "Token <easyverein-api-token>"
func SyncCSV(requestID string) {
	prefixID := fmt.Sprintf(prefix, requestID)

	file, err := os.Open("Mitglieder.csv")
	if err != nil {
		log.Println("sync failed, could no read csv file")
	}
	csvReader := csv.NewReader(file)
	csvReader.Comma = ';'

	members, err := csvReader.ReadAll()
	if err != nil {
		log.Println("sync failed, could not read csv")
	}

	// TODO implement csv reader and post to easy api
	log.Println(prefixID, members)
}

func makeAPIRequest(client *resty.Client, url string) (*resty.Response, error) {
	resp, err := client.R().
		SetHeader(
			"Authorization", fmt.Sprintf("Token %s", config.GetConfig().Easyverein.Token),
		).
		Get(url)
	if err != nil {
		return resp, err
	}

	return resp, nil
}
