package wordpress

import (
	"cmd/service/main.go/pkg/easyverein"
	"testing"
)

func TestGenerateLoginNames(t *testing.T) {
	testTable := []struct {
		testMember      easyverein.Member
		wantedLoginName string
	}{
		{
			easyverein.Member{LoginName: "", FirstName: "Max", LastName: "Mustermann", Email: ""},
			"max.mustermann",
		},
		{
			easyverein.Member{LoginName: "", FirstName: "Jürgen", LastName: "Ümläute", Email: ""},
			"juergen.uemlaeute",
		},
		{
			easyverein.Member{LoginName: "", FirstName: "Fred", LastName: "Feuerstein Meyer", Email: ""},
			"fred.feuerstein.meyer",
		},
		{
			easyverein.Member{LoginName: "", FirstName: "Don John", LastName: "McNamara", Email: ""},
			"don.john.mcnamara",
		},
		{
			easyverein.Member{LoginName: "", FirstName: "Bjorgen-Marie", LastName: "Kjörgen-Müller de Rapp", Email: ""},
			"bjorgen-marie.kjoergen-mueller.de.rapp",
		},
	}

	for _, table := range testTable {
		got := GenerateLoginName(table.testMember)

		if got != table.wantedLoginName {
			t.Errorf("generatedLoginName(%s %s) was incorrect, got: %s, want: %s",
				table.testMember.FirstName,
				table.testMember.LastName,
				got,
				table.wantedLoginName,
			)
		}
	}
}
