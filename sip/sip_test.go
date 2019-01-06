package sip

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/go-test/deep"
)

func TestInvitationCodec(t *testing.T) {
	msg := ControlMessage{
		Cmd:   Invitation,
		Token: 12345,
		SSRC:  54321,
		Name:  "foo-bar",
	}
	buffer := Encode(msg)
	actual, _ := Decode(buffer)

	fmt.Println(hex.Dump([]byte(msg.Name)))

	if diff := deep.Equal(msg, actual); diff != nil {
		t.Error(diff)
	}

}
