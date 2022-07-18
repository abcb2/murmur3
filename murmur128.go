package murmur3

import (
	"encoding/binary"
	"hash"
	"math/bits"
)

const (
	c1_128 = 0x87c37b91114253d5
	c2_128 = 0x4cf5ad432745937f
)

type sum128 struct {
	seed uint32
	h1   uint64
	h2   uint64
}

func New128() Hash128 {
	return New128WithSeed(0)
}

func New128WithSeed(seed uint32) Hash128 {
	s := sum128{
		seed: seed,
		h1:   0,
		h2:   0,
	}
	return &s
}

func (s *sum128) Reset() {
	s.h1, s.h2 = uint64(s.seed), uint64(s.seed)
}

func (s *sum128) Write(data []byte) (n int, err error) {
	dataLen := len(data)
	nBlocks := dataLen / 16
	start, end := 0, 0
	for i := 0; i < nBlocks; i++ {
		start = i * 16
		end = start + 16
		dataBlock := data[start:end]
		dataBlockLen := len(dataBlock)

		k1 := binary.LittleEndian.Uint64(dataBlock[0 : dataBlockLen/2])
		k2 := binary.LittleEndian.Uint64(dataBlock[dataBlockLen/2:])

		k1 *= c1_128
		k1 = bits.RotateLeft64(k1, 31)
		k1 *= c2_128
		s.h1 ^= k1

		s.h1 = bits.RotateLeft64(s.h1, 27)
		s.h1 += s.h2
		s.h1 = s.h1*5 + 0x52dce729

		k2 *= c2_128
		k2 = bits.RotateLeft64(k2, 33)
		k2 *= c1_128
		s.h2 ^= k2

		s.h2 = bits.RotateLeft64(s.h2, 31)
		s.h2 += s.h1
		s.h2 = s.h2*5 + 0x38495ab5
	}

	// tail
	tail := data[end:]
	var k1, k2 uint64 = 0, 0

	switch dataLen & 15 {
	case 15:
	case 14:
	case 13:
	case 12:
	case 11:
	}

	return dataLen, nil
}

func (s *sum128) Sum(b []byte) []byte {
	//TODO implement me
	panic("implement me")
}

func (s *sum128) Size() int {
	return 16
}

func (s *sum128) BlockSize() int {
	return 1
}

func (s *sum128) Sum128() (uint64, uint64) {
	return s.h1, s.h2
}

type Hash128 interface {
	hash.Hash
	Sum128() (uint64, uint64)
}
