package jacodoma

import (
	. "../src/"
	"github.com/smartystreets/assertions/should"
	. "github.com/smartystreets/goconvey/convey"
	"io/ioutil"
	"testing"
)

func TestParseParticipants(t *testing.T) {
	Convey("Parse line Leandro Santiago", t, func() {
		p, e := ParticipantFromString("Leandro Santiago <leandrosansilva@gmail.com>")
		So(e, should.Equal, nil)
		So(p.Name, should.Equal, "Leandro Santiago")
		So(p.Email, should.Equal, "leandrosansilva@gmail.com")
	})

	Convey("Ill formed: Parse line Leandro Santiago without space between name and e-mail", t, func() {
		p, e := ParticipantFromString("Leandro Santiago<leandrosansilva@gmail.com>")
		So(e, should.NotEqual, nil)
		So(p.Name, should.Equal, "")
		So(p.Email, should.Equal, "")
	})

	Convey("Ill former: no e-mail", t, func() {
		p, e := ParticipantFromString("Manoel Jobas")
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
		fileLines := `Leandro Santiago <leandrosansilva@gmail.com>
Manoel Ribas <manoel@ribas.go>
Juka Juke <juka@ju.ke>
Jon Doe <joe@doe.com>`

		err := ioutil.WriteFile("/tmp/list_with_4_participants.txt", []byte(fileLines), 0644)
		So(err, should.Equal, nil)
		participants, e := LoadParticipantsFromFile("/tmp/list_with_4_participants.txt")
		So(len(participants), should.Equal, 4)
		So(e, should.Equal, nil)
		So(participants[0].Name, should.Equal, "Leandro Santiago")
		So(participants[1].Name, should.Equal, "Manoel Ribas")
		So(participants[2].Name, should.Equal, "Juka Juke")
		So(participants[3].Name, should.Equal, "Jon Doe")
		So(participants[0].Email, should.Equal, "leandrosansilva@gmail.com")
		So(participants[1].Email, should.Equal, "manoel@ribas.go")
		So(participants[2].Email, should.Equal, "juka@ju.ke")
		So(participants[3].Email, should.Equal, "joe@doe.com")
	})
}
