package timestamp

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const tick = 100 * time.Microsecond

func Test_Of(t *testing.T) {
	// given
	start := time.Now()
	time := start.Add(tick)
	// when
	timestamp := Of(time, start)
	// then
	assert.Equal(t, uint64(1), timestamp.Uint64())
	assert.Equal(t, uint32(1), timestamp.Uint32())
}

func Test_Encode_one_octet_DeltaTime(t *testing.T) {
	// given
	b := new(bytes.Buffer)
	start := time.Now()
	reference := start.Add(tick)
	delta := tick

	// when
	EncodeDeltaTime(reference, start, delta, b)
	// then
	assert.Equal(t, []byte{0x01}, b.Bytes())
}

func Test_Encode_smallest_one_octet_DeltaTime(t *testing.T) {
	// given
	b := new(bytes.Buffer)
	start := time.Now()
	reference := start.Add(tick)
	delta := 99 * time.Microsecond

	// when
	EncodeDeltaTime(reference, start, delta, b)
	// then
	assert.Equal(t, []byte{0x00}, b.Bytes())
}

func Test_Encode_largest_one_octet_DeltaTime(t *testing.T) {
	// given
	b := new(bytes.Buffer)
	start := time.Now()
	reference := start.Add(tick)
	delta := 0x7f * tick

	// when
	EncodeDeltaTime(reference, start, delta, b)
	// then
	assert.Equal(t, []byte{0x7f}, b.Bytes())
}

func Test_Encode_smallest_two_octet_DeltaTime(t *testing.T) {
	// given
	b := new(bytes.Buffer)
	start := time.Now()
	reference := start.Add(tick)
	delta := 0x80 * tick

	// when
	EncodeDeltaTime(reference, start, delta, b)
	// then
	assert.Equal(t, []byte{0x81, 0x00}, b.Bytes())
}

func Test_Encode_largest_two_octet_DeltaTime(t *testing.T) {
	// given
	b := new(bytes.Buffer)
	start := time.Now()
	reference := start.Add(tick)
	delta := 0x3fff * tick

	// when
	EncodeDeltaTime(reference, start, delta, b)
	// then
	assert.Equal(t, []byte{0xff, 0x7f}, b.Bytes())
}

func Test_Encode_smallest_three_octet_DeltaTime(t *testing.T) {
	// given
	b := new(bytes.Buffer)
	start := time.Now()
	reference := start.Add(tick)
	delta := 0x4000 * tick

	// when
	EncodeDeltaTime(reference, start, delta, b)
	// then
	assert.Equal(t, []byte{0x81, 0x80, 0x00}, b.Bytes())
}

func Test_Encode_largest_three_octet_DeltaTime(t *testing.T) {
	// given
	b := new(bytes.Buffer)
	start := time.Now()
	reference := start.Add(tick)
	delta := 0x1fffff * tick

	// when
	EncodeDeltaTime(reference, start, delta, b)
	// then
	assert.Equal(t, []byte{0xff, 0xff, 0x7f}, b.Bytes())
}

func Test_Encode_smallest_four_octet_DeltaTime(t *testing.T) {
	// given
	b := new(bytes.Buffer)
	start := time.Now()
	reference := start.Add(tick)
	delta := 0x200000 * tick

	// when
	EncodeDeltaTime(reference, start, delta, b)
	// then
	assert.Equal(t, []byte{0x81, 0x80, 0x80, 0x00}, b.Bytes())
}

func Test_Encode_largest_three_four_DeltaTime(t *testing.T) {
	// given
	b := new(bytes.Buffer)
	start := time.Now()
	reference := start.Add(tick)
	delta := 0x0fffffff * tick

	// when
	EncodeDeltaTime(reference, start, delta, b)
	// then
	assert.Equal(t, []byte{0xff, 0xff, 0xff, 0x7f}, b.Bytes())
}
