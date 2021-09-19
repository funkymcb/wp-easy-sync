package models

import (
	"fmt"
	"log"
	"strings"
	"time"
)

// Member stores all necessarry data for wordpress account creation
type User struct {
	LoginName   string // will be generated. Default: firstname.lastname
	Password    string // will be generated. Default: ddmmyyyy
	FirstName   string `json:"firstName"`
	LastName    string `json:"familyName"`
	Email       string `json:"privateEmail,omitempty"`
	DateOfBirth string `json:"dateOfBirth,omitempty"`
}

// User.GenerateLoginName() of the convention: firstname.lastname
func (u User) GenerateLoginName() string {
	loginFirstName := replaceMutations(u.FirstName)
	loginLastName := replaceMutations(u.LastName)
	loginName := fmt.Sprintf("%s.%s",
		loginFirstName,
		loginLastName,
	)

	return loginName
}

// User.GeneratePassword() generated the default Password:
// Date of birth (ddmmyyyy)
func (u User) GeneratePassword() string {
	dateLayout := "2006-01-02"
	dateOfBirth, err := time.Parse(dateLayout, u.DateOfBirth)
	if err != nil {
		log.Println("Failed password creation. Date of Birth either missing or in the wrong format.")
		log.Printf("\tPassword of %s %s is set to 'default'", u.FirstName, u.LastName)
		return "default"
	}

	return dateOfBirth.Format("02.01.2006")
}

// replaces all known relevant special characters
// feel free to add more if needed
func replaceMutations(str string) string {
	str = strings.ToLower(str)
	str = strings.ReplaceAll(str, " ", ".")
	str = strings.ReplaceAll(str, "ä", "ae")
	str = strings.ReplaceAll(str, "ü", "ue")
	str = strings.ReplaceAll(str, "ö", "oe")
	str = strings.ReplaceAll(str, "ß", "ss")
	str = strings.ReplaceAll(str, "/", "-")
	str = strings.ReplaceAll(str, ".-.", "-")
	str = strings.ReplaceAll(str, "(", "-")
	str = strings.ReplaceAll(str, ")", "")
	str = strings.ReplaceAll(str, "dr..", "dr.")

	return str
}
