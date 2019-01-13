package timestamp

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_Of(t *testing.T) {
	// given
	start := time.Now()
	time := start.Add(100 * time.Microsecond)
	// when
	timestamp := Of(time, start)
	// then
	assert.Equal(t, uint64(1), timestamp.Uint64())
	assert.Equal(t, uint32(1), timestamp.Uint32())
}

func Test_EncodeDeltaTime(t *testing.T) {
	// given
	b := new(bytes.Buffer)
	start := time.Now()
	reference := start.Add(100 * time.Microsecond)
	delta := 100 * time.Microsecond

	// when
	EncodeDeltaTime(reference, start, delta, b)
	// then
	assert.Equal(t, []byte{0x01}, b.Bytes())
}
