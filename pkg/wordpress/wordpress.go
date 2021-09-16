package wordpress

import (
	"cmd/service/main.go/pkg/config"
	"cmd/service/main.go/pkg/models"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/go-resty/resty/v2"
)

type WordpressResponse []WPUser

type WPUser struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

var page = 1
var users []models.User

// GetUsers() unmarshals the API response of the wp users endpoint
// into a slice of Users
func GetUsers(client *resty.Client) ([]models.User, error) {
	var wordpressResponse WordpressResponse

	requestURI := config.GetConfig().Wordpress.APIRequestURI("users", page)

	resp, err := makeAPIRequest(client, requestURI)
	if err != nil {
		return users, fmt.Errorf("Could not perform GET Request to wp users endpoint: %v", err)
	}

	err = json.Unmarshal(resp.Body(), &wordpressResponse)
	if err != nil {
		return users, err
	}

	// total number of pages is given by the Response Header 'X-WP-TotalPages'
	numberOfPages, err := strconv.Atoi(resp.Header()["X-Wp-Totalpages"][0])
	if err != nil {
		return users, fmt.Errorf("Could not parse the Value of 'X-WP-Totalpages' header: %v", err)
	}

	// iterate from 2nd page to the last page (number of pages)
	if page < numberOfPages {
		for _, user := range wordpressResponse {
			users = append(users, models.User{
				LoginName: user.Username,
				FirstName: user.FirstName,
				LastName:  user.LastName,
				Email:     user.Email,
			})
		}
		page += 1
		GetUsers(client)
	} else if page == numberOfPages {
		//append members of the last page
		for _, user := range wordpressResponse {
			users = append(users, models.User{
				LoginName: user.Username,
				FirstName: user.FirstName,
				LastName:  user.LastName,
				Email:     user.Email,
			})
		}
	}

	return users, nil
}

// makes API GET request using resty
// headers
//   "Authorization": "Basic base64encoded(username:password)"
func makeAPIRequest(client *resty.Client, url string) (*resty.Response, error) {
	resp, err := client.R().
		SetBasicAuth(
			config.GetConfig().Wordpress.Username,
			config.GetConfig().Wordpress.Password,
		).
		Get(url)
	if err != nil {
		return resp, err
	}

	return resp, nil
}
