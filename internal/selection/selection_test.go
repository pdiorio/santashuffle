package selection

import (
	"math/rand"
	"testing"
)

func TestCreateParticipantsFromFile(t *testing.T) {
	participants, err := createParticipantsFromFile("./test.yaml")
	if err != nil {
		t.Errorf("createParticipantsFromFile() error = %v", err)
	}

	if len(participants) != 8 {
		t.Errorf("createParticipantsFromFile() got:%d, want:%d ", len(participants), 8)
	}
}

func TestCreateParticipantsFromYaml(t *testing.T) {
	s := `
participant 1:
    name: John Johnson
    email: john@example.com
    banned: ["Lucy Luciano"]
`
	participants, err := createParticipantsFromYaml([]byte(s))

	if err != nil {
		t.Errorf("createParticipantsFromYaml() error = %v", err)
	}

	if len(participants) != 1 {
		t.Errorf("createParticipantsFromYaml() got:%d, want:%d ", len(participants), 3)
	}

	if participants[0].Name != "John Johnson" {
		t.Errorf("createParticipantsFromYaml() got:%s, want:%s ", participants[0].Name, "John Johnson")
	}

	if participants[0].Email != "john@example.com" {
		t.Errorf("createParticipantsFromYaml() got:%s, want:%s ", participants[0].Email, "john@example.com")
	}

	if participants[0].Banned[0] != "Lucy Luciano" {
		t.Errorf("createParticipantsFromYaml() got:%s, want:%s ", participants[0].Banned[0], "Lucy Luciano")
	}

	if participants[0].identifier != "participant 1" {
		t.Errorf("createParticipantsFromYaml() got:%s, want:%s ", participants[0].identifier, "participant 1")
	}
}

func ExampleDisplayParticipants() {
	test := []*Participant{
		{Name: "Name1", Email: "Email1", Banned: []string{"BannedName"}},
		{Name: "Name2", Email: "Email2", Banned: []string{"BannedName"}},
		{Name: "Name3", Email: "Email3", Banned: []string{"BannedName"}},
		{Name: "Name4", Email: "Email4", Banned: []string{"BannedName"}},
	}

	DisplayParticipants(test)
	// Output:
	//0: Name1
	//1: Name2
	//2: Name3
	//3: Name4
}

func ExampleDisplayMatches() {
	test := []*Match{
		{Gifter: &Participant{Name: "Name1"}, Giftee: &Participant{Name: "Name3"}},
		{Gifter: &Participant{Name: "Name2"}, Giftee: &Participant{Name: "Name1"}},
		{Gifter: &Participant{Name: "Name3"}, Giftee: &Participant{Name: "Name2"}},
	}

	DisplayMatches(test)
	// Output:
	//Name1 -> Name3
	//Name2 -> Name1
	//Name3 -> Name2
}

func TestShuffleParticipants(t *testing.T) {
	rand.Seed(20)

	test := []*Participant{
		{Name: "Name1", Email: "Email1", Banned: []string{"BannedName"}},
		{Name: "Name2", Email: "Email2", Banned: []string{"BannedName"}},
		{Name: "Name3", Email: "Email3", Banned: []string{"BannedName"}},
	}

	shuffleParticipants(test)

	if (test[0].Name == "Name1") && (test[1].Name == "Name2") && (test[2].Name == "Name3") {
		t.Errorf("shuffleParticipants() participant list did not shuffle")
	}
}

func TestMakeSelections(t *testing.T) {
	rand.Seed(20)

	test := []*Participant{
		{Name: "Name1", Email: "Email1", Banned: []string{"BannedName"}},
		{Name: "Name2", Email: "Email2", Banned: []string{"BannedName"}},
		{Name: "Name3", Email: "Email3", Banned: []string{"BannedName"}},
	}
	matches := makeSelections(test)

	if len(matches) != 3 {
		t.Errorf("makeSelections() # of matches got:%d, want:%d ", len(matches), 3)
	}

	gifters := make(map[string]bool)
	giftees := make(map[string]bool)
	for _, match := range matches {
		gifters[match.Gifter.Name] = true
		giftees[match.Giftee.Name] = true
	}

	if len(gifters) != 3 {
		t.Errorf("makeSelections() # of gifters got:%d, want:%d ", len(gifters), 3)
	}
	if len(giftees) != 3 {
		t.Errorf("makeSelections() # of giftees got:%d, want:%d ", len(giftees), 3)
	}

}

func TestRandPositionWithExclusion(t *testing.T) {
	failed := false
	for i := 0; i < 2000; i++ {
		min, max, exclusion := 0, 20, 5
		position := randPositionWithExclusion(min, max, exclusion)
		if (position == exclusion) || (position < min) || (position > max) {
			failed = true
		}
	}

	if failed {
		t.Errorf("randPositionWithExclusion() failed bounds and exlusion check")
	}
}

func TestIsMatchAllowed(t *testing.T) {
	participantPtr1 := &Participant{Name: "Name1", Email: "Email1", Banned: []string{"Name3"}}
	participantPtr2 := &Participant{Name: "Name2", Email: "Email2"}
	participantPtr3 := &Participant{Name: "Name3", Email: "Email3", Banned: []string{"Name2"}}

	validMatch1 := &Match{Gifter: participantPtr1, Giftee: participantPtr2}
	validMatch2 := &Match{Gifter: participantPtr2, Giftee: participantPtr3}
	invalidMatch1 := &Match{Gifter: participantPtr1, Giftee: participantPtr3}
	invalidMatch2 := &Match{Gifter: participantPtr3, Giftee: participantPtr2}

	validResult1 := isMatchAllowed(validMatch1)
	validResult2 := isMatchAllowed(validMatch2)
	invalidResult1 := isMatchAllowed(invalidMatch1)
	invalidResult2 := isMatchAllowed(invalidMatch2)

	if !(validResult1 && validResult2) {
		t.Errorf("isMatchedAllowed() valid matches failed")
	}

	if invalidResult1 || invalidResult2 {
		t.Errorf("isMatchedAllowed() invalid matches pass")
	}

}

func TestSwapMatches(t *testing.T) {
	matchPtr1 := &Match{Gifter: &Participant{Name: "Name1"}, Giftee: &Participant{Name: "Name2"}}
	matchPtr2 := &Match{Gifter: &Participant{Name: "Name3"}, Giftee: &Participant{Name: "Name4"}}

	swapMatches(matchPtr1, matchPtr2)

	if !(matchPtr1.Gifter.Name == "Name1") || !(matchPtr2.Gifter.Name == "Name3") {
		t.Errorf("swapMatches() gifters swapped when only giftees should")
	}

	if !(matchPtr1.Giftee.Name == "Name4") || !(matchPtr2.Giftee.Name == "Name2") {
		t.Errorf("swapMatches() giftees did not swap correctly")
	}
}

func TestCorrectMatches(t *testing.T) {
	rand.Seed(20)
	test := []*Match{
		{Gifter: &Participant{Name: "Name1"}, Giftee: &Participant{Name: "Name1"}},
		{Gifter: &Participant{Name: "Name2"}, Giftee: &Participant{Name: "Name3"}},
		{Gifter: &Participant{Name: "Name3"}, Giftee: &Participant{Name: "Name4"}},
		{Gifter: &Participant{Name: "Name4", Banned: []string{"Name2"}}, Giftee: &Participant{Name: "Name2"}},
	}

	numCorrections, err := correctMatches(test)

	if err != nil {
		t.Errorf("correctMatches() error: %s", err)
	}

	if !(numCorrections >= 2) {
		t.Errorf("correctMatches() not enough corrections got:%d want>=%d", numCorrections, 2)
	}
}

func ExampleRunSelection() {
	_, _ = RunSelection("./test.yaml", 20, true)

	// Output:
	//Skylar White -> Garth Vader
	//Susan Luciano -> Lucy Johnson
	//Anne Vader -> John Johnson
	//Garth Vader -> Skylar White
	//Malcom White -> Roger Luciano
	//John Johnson -> Anne Vader
	//Roger Luciano -> Malcom White
	//Lucy Johnson -> Susan Luciano
	//total corrections made: 2
}
