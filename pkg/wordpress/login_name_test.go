package wordpress

import (
	"cmd/service/main.go/pkg/models"
	"testing"
)

func TestGenerateLoginNames(t *testing.T) {
	tests := []struct {
		testMember      models.User
		wantedLoginName string
	}{
		{
			models.User{LoginName: "", FirstName: "Max", LastName: "Mustermann", Email: ""},
			"max.mustermann",
		},
		{
			models.User{LoginName: "", FirstName: "Jürgen", LastName: "Ümläute", Email: ""},
			"juergen.uemlaeute",
		},
		{
			models.User{LoginName: "", FirstName: "Fred", LastName: "Feuerstein Meyer", Email: ""},
			"fred.feuerstein.meyer",
		},
		{
			models.User{LoginName: "", FirstName: "Don John", LastName: "McNamara", Email: ""},
			"don.john.mcnamara",
		},
		{
			models.User{LoginName: "", FirstName: "Bjorgen-Marie", LastName: "Kjörgen-Müller de Rapp", Email: ""},
			"bjorgen-marie.kjoergen-mueller.de.rapp",
		},
	}

	for _, test := range tests {
		got := GenerateLoginName(test.testMember)

		if got != test.wantedLoginName {
			t.Errorf("generatedLoginName(%s %s) was incorrect, got: %s, want: %s",
				test.testMember.FirstName,
				test.testMember.LastName,
				got,
				test.wantedLoginName,
			)
		}
	}
}
