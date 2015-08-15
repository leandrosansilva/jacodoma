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

func ParticipantFromString(s string) (Participant, error) {
	re := regexp.MustCompile("(.*)[[:space:]]+<(.*)>")
	matches := re.FindStringSubmatch(s)

	if len(matches) != 3 {
		return Participant{"", ""}, errors.New("Invalid Line")
	}

	return Participant{matches[1], matches[2]}, nil
}

func LoadParticipantsFromFile(filename string) ([]Participant, error) {
	file, e := os.Open(filename)

	if e != nil {
		return make([]Participant, 0), e
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	participants := make([]Participant, 0)

	for scanner.Scan() {
		s := scanner.Text()
		p, e := ParticipantFromString(s)
		if e != nil {
			return make([]Participant, 0), e
		}
		participants = append(participants, p)
	}

	return participants, nil
}
