package sip

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/go-test/deep"
	"github.com/stretchr/testify/assert"
)

func Test_Invitation_Codec(t *testing.T) {
	// given
	msg := ControlMessage{
		Cmd:   Invitation,
		Token: 12345,
		SSRC:  54321,
		Name:  "foo-bar",
	}
	// when
	buffer := Encode(msg)
	actual, err := Decode(buffer)
	// then
	fmt.Println(hex.Dump(buffer))
	assert.Nil(t, err)
	if diff := deep.Equal(msg, actual); diff != nil {
		t.Error(diff)
	}
}

func Test_Ignore_Name_In_End(t *testing.T) {
	// given
	msg := ControlMessage{
		Cmd:   End,
		Name:  "foo-bar",
	}
	// when
	buffer := Encode(msg)
	actual, err := Decode(buffer)
	// then
	fmt.Println(hex.Dump(buffer))
	assert.Nil(t, err)
	assert.Equal(t, actual.Name, "")
}

func Test_Timesync_Code(t *testing.T) {
	// given
	msg := ControlMessage{
		Cmd:   Synchronization,
		Timestamps: []uint64{0x0102030405060708, 0x1112131415161718 , 0x2122232425262728},
	}
	// when
	buffer := Encode(msg)
	actual, err := Decode(buffer)
	// then
	fmt.Println(hex.Dump(buffer))
	assert.Nil(t, err)
	if diff := deep.Equal(msg, actual); diff != nil {
		t.Error(diff)
	}
}
