package murmur3

import (
	"encoding/binary"
	"hash"
	"math/bits"
)

type (
	sum32 struct {
		seed uint32
		h    uint32
	}
)

func New32() hash.Hash32 {
	return New32WithSeed(0)
}

func New32WithSeed(seed uint32) hash.Hash32 {
	var s sum32 = sum32{
		seed: seed,
		h:    0,
	}
	return &s
}

func (s *sum32) Reset() {
	s.h = s.seed
}

func (s *sum32) Write(data []byte) (n int, err error) {
	s.h = s.seed

	dataLen := len(data)
	nBlocks := dataLen / 4
	start, end := 0, 0
	for i := 0; i < nBlocks; i++ {
		start = i * 4
		end = start + 4

		k := binary.LittleEndian.Uint32(data[start:end])

		k = k * 0xcc9e2d51
		k = bits.RotateLeft32(k, 15)
		k = k * 0x1b873593

		s.h = s.h ^ k
		s.h = bits.RotateLeft32(s.h, 13)
		s.h = s.h*5 + 0xe6546b64
	}

	// tail
	tail := data[end:]
	var k1 uint32 = 0
	switch dataLen & 3 {
	case 3:
		k1 ^= uint32(tail[2]) << 16
		fallthrough
	case 2:
		k1 ^= uint32(tail[1]) << 8
		fallthrough
	case 1:
		k1 ^= uint32(tail[0])
		k1 *= 0xcc9e2d51
		k1 = bits.RotateLeft32(k1, 15)
		k1 *= 0x1b873593
		s.h ^= k1
	}

	// finalize
	s.h ^= uint32(dataLen)
	s.h ^= s.h >> 16
	s.h *= 0x85ebca6b
	s.h ^= s.h >> 13
	s.h *= 0xc2b2ae35
	s.h ^= s.h >> 16

	return dataLen, nil
}

func (s *sum32) Sum(in []byte) []byte {
	return append(in, byte(s.h>>24), byte(s.h>>16), byte(s.h>>8), byte(s.h))
}

func (s *sum32) Size() int {
	return 4
}

func (s *sum32) BlockSize() int {
	return 1
}

func (s *sum32) Sum32() uint32 {
	return s.h
}
