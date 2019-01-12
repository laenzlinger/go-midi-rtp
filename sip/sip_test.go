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
	buffer, err := Encode(msg)
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
		Cmd:  End,
		Name: "foo-bar",
	}
	// when
	buffer, err := Encode(msg)
	actual, err := Decode(buffer)
	// then
	fmt.Println(hex.Dump(buffer))
	assert.Nil(t, err)
	assert.Equal(t, actual.Name, "")
}

func Test_Timesync_Codec(t *testing.T) {
	// given
	msg := ControlMessage{
		Cmd:        Synchronization,
		SSRC:       0xaabbccdd,
		Timestamps: []uint64{0x0102030405060708, 0x1112131415161718, 0x2122232425262728},
	}
	// when
	buffer, err := Encode(msg)
	actual, err := Decode(buffer)
	// then
	fmt.Println(hex.Dump(buffer))
	assert.Nil(t, err)
	assert.Equal(t, len(buffer), 36)
	if diff := deep.Equal(msg, actual); diff != nil {
		t.Error(diff)
	}
}

func Test_Timesync_Encoding_Should_Send_Complete_Package(t *testing.T) {
	// given
	msg := ControlMessage{
		Cmd:        Synchronization,
		Timestamps: []uint64{0x1111111111111111},
	}
	// when
	buffer, err := Encode(msg)
	actual, err := Decode(buffer)
	// then
	fmt.Println(hex.Dump(buffer))
	assert.Nil(t, err)
	if diff := deep.Equal(msg, actual); diff != nil {
		t.Error(diff)
	}
	assert.Equal(t, []byte{
		0xff, 0xff, 0x43, 0x4b, // header | cmd
		0x00, 0x00, 0x00, 0x00, // SSRC
		0x00, 0x00, 0x00, 0x00, // count
		0x11, 0x11, 0x11, 0x11, // timstamp 1 (high)
		0x11, 0x11, 0x11, 0x11, // timstamp 1 (low)
		0x00, 0x00, 0x00, 0x00, // timstamp 2 (high)
		0x00, 0x00, 0x00, 0x00, // timstamp 2 (low)
		0x00, 0x00, 0x00, 0x00, // timstamp 3 (high)
		0x00, 0x00, 0x00, 0x00, // timstamp 3 (low)
	}, buffer)
	assert.Equal(t, 36, len(buffer))
}

func Test_Timesync_Encoding_Without_Timestamp_is_wrong(t *testing.T) {
	// given
	msg := ControlMessage{
		Cmd: Synchronization,
	}
	// when
	_, err := Encode(msg)
	// then
	assert.Error(t, err)
}
