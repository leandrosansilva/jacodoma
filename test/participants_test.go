package jacodoma

import (
	. "../src/"
	"github.com/smartystreets/assertions/should"
	. "github.com/smartystreets/goconvey/convey"
	"io/ioutil"
	"testing"
)

func TestParseParticipants(t *testing.T) {
	Convey("Parse line Coding Dojo", t, func() {
		p, e := ParticipantFromString("Coding Dojo <coding@do.jo>")
		So(e, should.Equal, nil)
		So(p.Name, should.Equal, "Coding Dojo")
		So(p.Email, should.Equal, "coding@do.jo")
	})

	Convey("Ill formed: Parse line Coding Dojo without space between name and e-mail", t, func() {
		p, e := ParticipantFromString("Coding Dojo<coding@do.jo>")
		So(e, should.NotEqual, nil)
		So(p.Name, should.Equal, "")
		So(p.Email, should.Equal, "")
	})

	Convey("Ill former: no e-mail", t, func() {
		p, e := ParticipantFromString("Juca Pinto")
		So(e, should.NotEqual, nil)
		So(p.Name, should.Equal, "")
		So(p.Email, should.Equal, "")
	})

	Convey("Ill former: no name", t, func() {
		p, e := ParticipantFromString("<juca@pinto.com>")
		So(e, should.NotEqual, nil)
		So(p.Name, should.Equal, "")
		So(p.Email, should.Equal, "")
	})
}

func TestReadParticipantListFromFile(t *testing.T) {
	Convey("Parse Empty File", t, func() {
		// FIXME: this won't work on non-unixes systems!
		participants, e := LoadParticipantsFromFile("/dev/null")
		So(len(participants), should.Equal, 0)
		So(e, should.Equal, nil)
	})

	Convey("Parse File with 4 participants", t, func() {
		fileLines := `Coding Dojo <coding@do.jo>
Manoel Ribas <manoel@ribas.go>
Juka Juke <juka@ju.ke>
Jon Doe <joe@doe.com>`

		err := ioutil.WriteFile("/tmp/list_with_4_participants.txt", []byte(fileLines), 0644)
		So(err, should.Equal, nil)
		participants, e := LoadParticipantsFromFile("/tmp/list_with_4_participants.txt")
		So(len(participants), should.Equal, 4)
		So(e, should.Equal, nil)
		So(participants[0].Name, should.Equal, "Coding Dojo")
		So(participants[1].Name, should.Equal, "Manoel Ribas")
		So(participants[2].Name, should.Equal, "Juka Juke")
		So(participants[3].Name, should.Equal, "Jon Doe")
		So(participants[0].Email, should.Equal, "coding@do.jo")
		So(participants[1].Email, should.Equal, "manoel@ribas.go")
		So(participants[2].Email, should.Equal, "juka@ju.ke")
		So(participants[3].Email, should.Equal, "joe@doe.com")
	})

	Convey("Error parsing file", t, func() {
		fileLines := `Coding Dojo <coding@do.jo>
Manoel Ribas<manoel@ribas.go>
Jon Doe <joe@doe.com>`

		err := ioutil.WriteFile("/tmp/list_with_3_participants_and_error.txt", []byte(fileLines), 0644)
		So(err, should.Equal, nil)
		participants, e := LoadParticipantsFromFile("/tmp/list_with_3_participants_and_error")
		So(len(participants), should.Equal, 0)
		So(e, should.NotEqual, nil)
	})

}
