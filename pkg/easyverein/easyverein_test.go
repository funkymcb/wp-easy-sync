package easyverein

import "testing"

func TestGenerateLoginNames(t *testing.T) {
	testTable := []struct {
		testMember      Member
		wantedLoginName string
	}{
		{
			Member{"", "Max", "Mustermann", ""},
			"max.mustermann",
		},
		{
			Member{"", "Jürgen", "Ümläute", ""},
			"juergen.uemlaeute",
		},
		{
			Member{"", "Fred", "Feuerstein Meyer", ""},
			"fred.feuerstein.meyer",
		},
		{
			Member{"", "Don John", "McNamara", ""},
			"don.john.mcnamara",
		},
		{
			Member{"", "Bjorgen-Marie", "Kjörgen-Müller de Rapp", ""},
			"bjorgen-marie.kjoergen-mueller.de.rapp",
		},
	}

	for _, table := range testTable {
		got := generateLoginName(table.testMember)

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
