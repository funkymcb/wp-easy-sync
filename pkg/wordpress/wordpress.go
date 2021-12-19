package wordpress

import (
	"cmd/service/main.go/pkg/config"
	"cmd/service/main.go/pkg/models"
	"encoding/json"
	"fmt"
	"log"
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

var Page = 1
var users *[]models.User

func init() {
	users = &[]models.User{}
}

// GetUsers() unmarshals the API response of the wp users endpoint
// into a slice of Users
func GetUsers(client *resty.Client) (*[]models.User, error) {
	log.Printf("Fetching users from page: %d", Page)
	var wordpressResponse WordpressResponse

	requestURI := config.GetConfig().Wordpress.APIGETRequestURI("users", Page)

	resp, err := makeGETRequest(client, requestURI)
	if err != nil {
		return users, fmt.Errorf("could not perform get request to wp users endpoint: %v", err)
	}

	if resp.StatusCode() >= 400 {
		err = fmt.Errorf("status code: %d response body: %s", resp.StatusCode(), string(resp.Body()))
		return users, err
	}

	err = json.Unmarshal(resp.Body(), &wordpressResponse)
	if err != nil {
		return users, err
	}

	// total number of pages is given by the Response Header 'X-WP-TotalPages'
	numberOfPages, err := strconv.Atoi(resp.Header()["X-Wp-Totalpages"][0])
	if err != nil {
		return users, fmt.Errorf("could not parse the value of 'x-wp-totalpages' header: %v", err)
	}

	// iterate from 2nd page to the last page (number of pages)
	if Page < numberOfPages {
		for _, user := range wordpressResponse {
			*users = append(*users, models.User{
				LoginName: user.Username,
				FirstName: user.FirstName,
				LastName:  user.LastName,
				Email:     user.Email,
			})
		}
		Page += 1
		GetUsers(client)
	} else if Page == numberOfPages {
		//append members of the last page
		for _, user := range wordpressResponse {
			*users = append(*users, models.User{
				LoginName: user.Username,
				FirstName: user.FirstName,
				LastName:  user.LastName,
				Email:     user.Email,
			})
		}
	}

	return users, nil
}

// CreateUser() creates a Wordpress User Account
func CreateUser(client *resty.Client, user models.User) error {
	requestURI := config.GetConfig().Wordpress.APIPOSTRequestURI("users")

	postData := map[string]string{
		"username":   user.LoginName,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"email":      user.Email,
		"password":   user.Password,
	}

	resp, err := makePOSTRequest(client, requestURI, postData)
	if err != nil {
		return fmt.Errorf("Error performing API POST Request %v", err)
	}

	if resp.StatusCode() == 200 {
		log.Printf("Account creation for %s successful", user.LoginName)
	} else if resp.StatusCode() >= 400 && resp.StatusCode() < 500 {
		return fmt.Errorf("Account creation Request failed: Status Code: %d",
			resp.StatusCode(),
		)
	} else if resp.StatusCode() == 500 {
		log.Printf("Error creating User %s\nSkipping Account creation.\n Error: %s",
			user.LoginName,
			resp,
		)
	}

	return nil
}

// makes API GET request using resty
// headers
//   "Authorization": "Basic base64encoded(username:password)"
func makeGETRequest(client *resty.Client, url string) (*resty.Response, error) {
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

// makes API POST request using resty
// headers
//   "Authorization": "Basic base64encoded(username:password)"
// body
//	 string
func makePOSTRequest(client *resty.Client, url string, postData map[string]string) (*resty.Response, error) {
	resp, err := client.R().
		SetBasicAuth(
			config.GetConfig().Wordpress.Username,
			config.GetConfig().Wordpress.Password,
		).
		SetMultipartFormData(postData).
		Post(url)
	if err != nil {
		return resp, err
	}

	return resp, nil
}
