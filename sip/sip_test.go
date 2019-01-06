package sip

import (
	"testing"

	"github.com/go-test/deep"
)

func TestMarshall(t *testing.T) {
	msg := ControlMessage{
		Cmd:   Invitation,
		Token: 12345,
		SSRC:  54321,
		Name:  "foo-bar",
	}
	buffer := Marshall(msg)
	actual, _ := Parse(buffer)

	if diff := deep.Equal(msg, actual); diff != nil {
		t.Error(diff)
	}

}
