package jacodoma

import (
	"bufio"
	"errors"
	"math/rand"
	"os"
	"regexp"
)

type Participant struct {
	Name  string
	Email string
}

type Participants struct {
	participants []Participant
}

func (P *Participants) Length() int {
	return len(P.participants)
}

func (P *Participants) Get(index int) Participant {
	return P.participants[index]
}

func (P *Participant) Valid() bool {
	return len(P.Name) > 0 && len(P.Email) > 0
}

func ParticipantFromString(s string) (Participant, error) {
	re := regexp.MustCompile("^[[:space:]]*(.*[^[:space:]])[[:space:]]*<(.*)>")
	matches := re.FindStringSubmatch(s)

	if len(matches) != 3 {
		return Participant{"", ""}, errors.New("Ill-formed line: \"" + s + "\"")
	}

	return Participant{matches[1], matches[2]}, nil
}

func LoadParticipantsFromFile(filename string) (Participants, error) {
	file, e := os.Open(filename)

	var empty Participants

	if e != nil {
		return empty, e
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	participants := make([]Participant, 0)

	for scanner.Scan() {
		p, e := ParticipantFromString(scanner.Text())

		if e != nil {
			return empty, e
		}

		participants = append(participants, p)
	}

	return Participants{participants}, nil
}

func (this *Participants) Shuffle() {
	for i := range this.participants {
		j := rand.Intn(i + 1)
		this.participants[i], this.participants[j] = this.participants[j], this.participants[i]
	}
}

func BuildParticipantsFromArray(participants []Participant) Participants {
	return Participants{participants}
}
