package main

import (
	. "./src"
	"fmt"
)

func init() {
}

func main() {
	participants, _ := LoadParticipantsFromFile("users.jcdm")

	for i := 1; i < participants.Length(); i++ {
		fmt.Println(participants.Get(i))
	}
}
