package models

import (
	"testing"
)

func TestGenerateLoginNames(t *testing.T) {
	tests := []struct {
		testMember      User
		wantedLoginName string
	}{
		{
			User{LoginName: "", FirstName: "Max", LastName: "Mustermann", Email: ""},
			"max.mustermann",
		},
		{
			User{LoginName: "", FirstName: "Jürgen", LastName: "Ümläute", Email: ""},
			"juergen.uemlaeute",
		},
		{
			User{LoginName: "", FirstName: "Fred", LastName: "Feuerstein Meyer", Email: ""},
			"fred.feuerstein.meyer",
		},
		{
			User{LoginName: "", FirstName: "Don John", LastName: "McNamara", Email: ""},
			"don.john.mcnamara",
		},
		{
			User{LoginName: "", FirstName: "Bjorgen-Marie", LastName: "Kjörgen-Müller de Rapp", Email: ""},
			"bjorgen-marie.kjoergen-mueller.de.rapp",
		},
	}

	for _, test := range tests {
		got := test.testMember.GenerateLoginName()

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
