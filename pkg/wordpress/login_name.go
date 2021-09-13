package wordpress

import (
	"cmd/service/main.go/pkg/models"
	"fmt"
	"strings"
)

// GenerateLoginName of the convention: firstname.lastname
func GenerateLoginName(member models.WVCMember) string {
	loginFirstName := replaceMutations(member.FirstName)
	loginLastName := replaceMutations(member.LastName)
	loginName := fmt.Sprintf("%s.%s",
		loginFirstName,
		loginLastName,
	)

	return loginName
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

	return str
}
