package easyverein

import (
	"cmd/service/main.go/pkg/config"
	"encoding/json"
	"fmt"

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

var page = 1
var members []Member

// ListMembers() unmarshals the API response of the contact-details endpoint
// into a slice of Members
func ListMembers(client *resty.Client) ([]Member, error) {
	var easyResponse EasyVereinResponse

	// requestURI = https://easyverein.com/api/stable/contact-details?limit100&page=%d
	requestURI := config.GetConfig().Easyverein.APIRequestURI("contact-details", page)

	resp, err := makeAPIRequest(client, requestURI)

	err = json.Unmarshal(resp.Body(), &easyResponse)
	if err != nil {
		return members, err
	}

	// call ListMembers() recursively until no next page
	if easyResponse.Next != "" {
		members = append(members, easyResponse.Members...)
		page += 1
		ListMembers(client)
	} else {
		// append members of the last page
		members = append(members, easyResponse.Members...)
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
