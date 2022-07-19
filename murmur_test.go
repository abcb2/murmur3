package murmur3

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

var data = []struct {
	seed  uint32
	h32   uint32
	h64_1 uint64
	h64_2 uint64
	s     string
}{
	{0x00, 0x00000000, 0x0000000000000000, 0x0000000000000000, ""},
	{0x00, 0x248bfa47, 0xcbd8a7b341bd9b02, 0x5b1e906a48ae1d19, "hello"},
}

func TestMurmur3(t *testing.T) {
	for idx, elem := range data {
		t.Run("h32", func(t *testing.T) {
			h32 := New32WithSeed(elem.seed)
			_, err := h32.Write([]byte(elem.s))
			assert.Nil(t, err)

			hashed := h32.Sum(nil)
			expectHexStr := fmt.Sprintf("%08x", elem.h32)
			resultHexStr := fmt.Sprintf("%08x", hashed)
			assert.Equal(t, expectHexStr, resultHexStr, idx)
		})
		t.Run("h128", func(t *testing.T) {
			h128 := New128WithSeed(elem.seed)
			_, _ = h128.Write([]byte(elem.s))
			v1, v2 := h128.Sum128()
			assert.Equal(t, elem.h64_1, v1)
			assert.Equal(t, elem.h64_2, v2)
		})
	}
}
