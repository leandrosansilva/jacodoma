package jacodoma

import (
	"bufio"
	"errors"
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

func ParticipantFromString(s string) (Participant, error) {
	re := regexp.MustCompile("(.*) <(.*)>")
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
