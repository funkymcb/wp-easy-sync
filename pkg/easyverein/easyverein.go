package easyverein

import (
	"cmd/service/main.go/pkg/config"
	"cmd/service/main.go/pkg/models"
	"encoding/json"
	"fmt"
	"log"

	"github.com/go-resty/resty/v2"
)

// EasyVereinResponse stores all necessarry values of API response
type EasyVereinResponse struct {
	// the url for the next page, null if last page
	Next    string        `json:"next"`
	Members []models.User `json:"results"`
}

var page = 1
var members []models.User

// GetMembers() unmarshals the API response of the contact-details endpoint
// into a slice of Users
func GetMembers(client *resty.Client) ([]models.User, error) {
	log.Printf("Fetching users from page: %d", page)
	var easyResponse EasyVereinResponse

	// requestURI = https://easyverein.com/api/stable/contact-details?limit100&page=%d
	requestURI := config.GetConfig().Easyverein.APIRequestURI("contact-details", page)

	resp, err := makeAPIRequest(client, requestURI)
	if err != nil {
		return members, fmt.Errorf("could not perform get request to easyverein contact-details endpoint: %v", err)
	}

	err = json.Unmarshal(resp.Body(), &easyResponse)
	if err != nil {
		return members, err
	}

	// call GetMembers() recursively until no next page
	if easyResponse.Next != "" {
		members = append(members, easyResponse.Members...)
		page += 1
		GetMembers(client)
	} else {
		// append members of the last page
		members = append(members, easyResponse.Members...)
	}

	return members, nil
}

// makes API GET request using resty
// headers:
//   "Authorization": "Token <easyverein-api-token>"
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
