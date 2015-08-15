package main

import (
	. "./src"
	"fmt"
)

func main() {
	participants, _ := LoadParticipantsFromFile("users.jcdm")

	for _, p := range participants {
		fmt.Println(p)
	}
}
