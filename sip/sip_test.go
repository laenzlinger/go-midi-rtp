package sip

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/go-test/deep"
	"github.com/stretchr/testify/assert"
)

func TestInvitationCodec(t *testing.T) {
	// given
	msg := ControlMessage{
		Cmd:   Invitation,
		Token: 12345,
		SSRC:  54321,
		Name:  "foo-bar",
	}
	// when
	actual, err := Decode(Encode(msg))
	// then
	fmt.Println(hex.Dump([]byte(msg.Name)))
	assert.Nil(t, err)
	if diff := deep.Equal(msg, actual); diff != nil {
		t.Error(diff)
	}
}
