package wordpress

import (
	"cmd/service/main.go/pkg/models"

	"github.com/go-resty/resty/v2"
)

//TODO analyze wordpress response for necessary fields:
type WordpressResponse struct {
}

func GetUsers(client *resty.Client) ([]models.WordpressUser, error) {
	var users []models.WordpressUser

	return users, nil
}
