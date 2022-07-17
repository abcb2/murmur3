package murmur3

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

var data = []struct {
	seed uint32
	h32  uint32
	s    string
}{
	{0x00, 0x00000000, ""},
	{0x00, 0x248bfa47, "hello"},
}

func TestMurmur3(t *testing.T) {
	for idx, elem := range data {
		h32 := New32WithSeed(elem.seed)
		_, err := h32.Write([]byte(elem.s))
		assert.Nil(t, err)

		hashed := h32.Sum(nil)
		expectHexStr := fmt.Sprintf("%08x", elem.h32)
		resultHexStr := fmt.Sprintf("%08x", hashed)
		assert.Equal(t, expectHexStr, resultHexStr, idx)
	}
}
