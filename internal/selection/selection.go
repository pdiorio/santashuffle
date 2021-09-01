package selection

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"sort"
	"time"

	"gopkg.in/yaml.v3"
)

type Participant struct {
	Name       string
	Email      string
	Banned     []string
	identifier string
}

type Match struct {
	Gifter *Participant
	Giftee *Participant
}

func createParticipantsFromFile(filename string) ([]*Participant, error) {
	yamlData, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return createParticipantsFromYaml(yamlData)
}

func createParticipantsFromYaml(yfile []byte) ([]*Participant, error) {
	data := make(map[string]Participant)
	err := yaml.Unmarshal(yfile, &data)

	if err != nil {
		return nil, err
	}

	keys := make([]string, 0)
	for k, _ := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	participants := make([]*Participant, 0, len(data))
	for _, k := range keys {
		participant := &Participant{
			Name:       data[k].Name,
			Email:      data[k].Email,
			Banned:     data[k].Banned,
			identifier: k,
		}
		participants = append(participants, participant)
	}

	return participants, nil

}

func DisplayParticipants(participants []*Participant) {
	for i, participant := range participants {
		fmt.Printf("%d: %s\n", i, participant.Name)
	}
}

func DisplayMatches(matches []*Match) {
	for _, match := range matches {
		fmt.Printf("%s -> %s\n", match.Gifter.Name, match.Giftee.Name)
	}
}

func shuffleParticipants(participants []*Participant) {
	rand.Shuffle(len(participants), func(i, j int) {
		participants[i], participants[j] = participants[j], participants[i]
	})
}

func makeSelections(participants []*Participant) []*Match {
	cpy := make([]*Participant, len(participants))
	copy(cpy, participants)

	shuffleParticipants(participants)
	shuffleParticipants(cpy)

	selections := make([]*Match, 0, len(participants))

	for i := range participants {
		match := &Match{
			Gifter: participants[i],
			Giftee: cpy[i],
		}
		selections = append(selections, match)
	}
	return selections
}

func randPositionWithExclusion(min int, max int, exclude int) int {
	i := exclude

	for i == exclude {
		i = rand.Intn(max-min) + min
	}

	return i
}

func isMatchAllowed(match *Match) bool {
	if match.Gifter.Name == match.Giftee.Name {
		return false
	}

	for _, bannedParticipant := range match.Gifter.Banned {
		if bannedParticipant == match.Giftee.Name {
			return false
		}
	}

	return true
}

func swapMatches(match1 *Match, match2 *Match) {
	match1.Giftee, match2.Giftee = match2.Giftee, match1.Giftee
}

func correctMatches(matches []*Match) int {
	reviewNeeded := true
	totalCorrections := 0

	for reviewNeeded {
		reviewNeeded = false
		for i, match := range matches {
			if !isMatchAllowed(match) {
				swapPosition := randPositionWithExclusion(0, len(matches)-1, i)
				swapMatches(match, matches[swapPosition])

				totalCorrections = totalCorrections + 1
				reviewNeeded = true
			}
		}
	}
	return totalCorrections
}

func RunSelection(participantFile string, seed int64, dryrun bool) []*Match {
	if seed == 0 {
		seed = time.Now().UnixNano()
	}
	rand.Seed(seed)

	participants, _ := createParticipantsFromFile(participantFile)

	matches := makeSelections(participants)
	totalCorrections := correctMatches(matches)

	if dryrun {
		DisplayMatches(matches)
		fmt.Printf("total corrections made: %d\n", totalCorrections)
	}

	return matches
}
