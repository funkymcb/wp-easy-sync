package easyverein

import (
	"cmd/service/main.go/pkg/config"
	"encoding/json"
	"fmt"
	"log"

	"github.com/go-resty/resty/v2"
)

// EasyVereinResponse stores all necessarry values of API response
type EasyVereinResponse struct {
	// the url for the next page, null if last page
	Next    string   `json:"next"`
	Members []Member `json:"results"`
}

// Member stores all necessarry data for wordpress account creation
type Member struct {
	FirstName  string `json:"firstName"`
	FamilyName string `json:"familyName"`
	Email      string `json:"privateEmail,omitempty"`
}

// ListMembers() unmarshals the API response of the contact-details endpoint
// into a slice of Members
func ListMembers() (EasyVereinResponse, error) {
	var members EasyVereinResponse

	// requestURI = https://easyverein.com/api/stable/contact-details?limit650
	requestURI := config.GetConfig().Easyverein.APIRequestURI("contact-details")

	client := resty.New()

	resp, err := makeAPIRequest(client, requestURI)

	err = json.Unmarshal(resp.Body(), &members)
	if err != nil {
		return members, err
	}

	//TODO fix iterating over all pages
	// current behaviour: infite loop of requesting page 2

	// keep requesting api and append members until "next" is null
	// in other words: iterate over all pages until we are on the last page
	for members.Next != "null" {
		resp, err = makeAPIRequest(client, members.Next)

		log.Println(resp)
		log.Println()
		log.Println()
	}

	return members, nil
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
