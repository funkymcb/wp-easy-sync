package easyverein

import (
	"cmd/service/main.go/pkg/config"
	"encoding/json"
	"fmt"
	"strings"

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
	LoginName string
	FirstName string `json:"firstName"`
	LastName  string `json:"familyName"`
	Email     string `json:"privateEmail,omitempty"`
}

var page = 1
var members []Member

// GetMembers() unmarshals the API response of the contact-details endpoint
// into a slice of Members
func GetMembers(client *resty.Client) ([]Member, error) {
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
		GetMembers(client)
	} else {
		// append members of the last page
		members = append(members, easyResponse.Members...)
	}

	for _, member := range members {
		member.LoginName = generateLoginName(member)
	}

	return members, nil
}

// makes API request using resty
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

// generates login name of the convention: firstname.lastname
func generateLoginName(member Member) string {
	loginFirstName := replaceMutations(member.FirstName)
	loginLastName := replaceMutations(member.LastName)
	loginName := fmt.Sprintf("%s.%s",
		loginFirstName,
		loginLastName,
	)

	return loginName
}

func replaceMutations(str string) string {
	str = strings.ToLower(str)
	str = strings.ReplaceAll(str, " ", ".")
	str = strings.ReplaceAll(str, "ä", "ae")
	str = strings.ReplaceAll(str, "ü", "ue")
	str = strings.ReplaceAll(str, "ö", "oe")
	str = strings.ReplaceAll(str, "ß", "ss")

	return str
}
