package main

import (
	. "./src"
	"fmt"
)

func main() {
	participants, _ := LoadParticipantsFromFile("users.jcdm")

	for p, _ := range participants {
		fmt.Println(p)
	}
}
