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
		k2 ^= uint64(tail[14]) << 48
		fallthrough
	case 14:
		k2 ^= uint64(tail[13]) << 40
		fallthrough
	case 13:
		k2 ^= uint64(tail[12]) << 32
		fallthrough
	case 12:
		k2 ^= uint64(tail[11]) << 24
		fallthrough
	case 11:
		k2 ^= uint64(tail[10]) << 16
		fallthrough
	case 10:
		k2 ^= uint64(tail[9]) << 8
		fallthrough
	case 9:
		k2 ^= uint64(tail[8]) << 0

		k2 *= c2_128
		k2 = bits.RotateLeft64(k2, 33)
		k2 *= c1_128
		s.h2 ^= k2
		fallthrough
	case 8:
		k1 ^= uint64(tail[7]) << 56
		fallthrough
	case 7:
		k1 ^= uint64(tail[6]) << 48
		fallthrough
	case 6:
		k1 ^= uint64(tail[5]) << 40
		fallthrough
	case 5:
		k1 ^= uint64(tail[4]) << 32
		fallthrough
	case 4:
		k1 ^= uint64(tail[3]) << 24
		fallthrough
	case 3:
		k1 ^= uint64(tail[2]) << 16
		fallthrough
	case 2:
		k1 ^= uint64(tail[1]) << 8
		fallthrough
	case 1:
		k1 ^= uint64(tail[0]) << 0
		k1 *= c1_128
		k1 = bits.RotateLeft64(k1, 31)
		k1 *= c2_128
		s.h1 ^= k1
	}

	s.h1 ^= uint64(dataLen)
	s.h2 ^= uint64(dataLen)

	s.h1 += s.h2
	s.h2 += s.h1

	s.h1 = fmix64(s.h1)
	s.h2 = fmix64(s.h2)

	s.h1 += s.h2
	s.h2 += s.h1

	return dataLen, nil
}

func fmix64(k uint64) uint64 {
	k ^= k >> 33
	k *= 0xff51afd7ed558ccd
	k ^= k >> 33
	k *= 0xc4ceb9fe1a85ec53
	k ^= k >> 33
	return k
}

func (s *sum128) Sum(b []byte) []byte {
	h1, h2 := s.h1, s.h2
	return append(b,
		byte(h1>>56), byte(h1>>48), byte(h1>>40), byte(h1>>32),
		byte(h1>>24), byte(h1>>16), byte(h1>>8), byte(h1),

		byte(h2>>56), byte(h2>>48), byte(h2>>40), byte(h2>>32),
		byte(h2>>24), byte(h2>>16), byte(h2>>8), byte(h2),
	)
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
