// Copyright (c) 2015-2016 The Decred developers
// Copyright (c) 2016 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package chainhash

import (
	"crypto/sha256"
	"github.com/dchest/blake256"
)

// HashFunc calculates the hash of the supplied bytes.
// TODO(jcv) Should modify blake256 so it has the same interface as blake2
// and sha256 so these function can look more like btcsuite.  Then should
// try to get it to the upstream blake256 repo
func HashFunc(b []byte) [blake256.Size]byte {
	var outB [blake256.Size]byte
	copy(outB[:], HashB(b))

	return outB
}

// HashB calculates hash(b) and returns the resulting bytes.
func HashB(b []byte) []byte {
	a := blake256.New()
	a.Write(b)
	out := a.Sum(nil)
	return out
}

// HashH calculates hash(b) and returns the resulting bytes as a Hash.
func HashH(b []byte) Hash {
	return Hash(HashFunc(b))
}

// XORBytes modifies b1 by pefroming a bitwise XOR with respective elements in b2
func XORBytes(b1 []byte, b2 []byte) {
	var size int

	if len(b1) < len(b2) {
		size = len(b1)
	} else {
		size = len(b2)
	}

	for i := 0; i < size; i++ {
		b1[i] = b1[i] ^ b2[i]
	}
}

// PoWHashB calculates (hash(hash(b) ^ MagicBytes) ^ MagicBytes) and returns the resulting bytes
func PoWHashB(b []byte) []byte {
	first := sha256.Sum256(b)
	XORBytes(first[:], MagicBytes[:])

	second := sha256.Sum256(first[:])
	XORBytes(second[:], MagicBytes[:])

	return second[:]
}

// PoWHashH calculates (hash(hash(b) ^ MagicBytes) ^ MagicBytes) and returns the resulting bytes as a
// Hash.
func PoWHashH(b []byte) Hash {
	var array [HashSize]byte
	copy(array[:], PoWHashB(b[:]))

	return Hash(array)
}

// HashBlockSize is the block size of the hash algorithm in bytes.
const HashBlockSize = blake256.BlockSize

// MagicBytes is the magic 256-bit random number used for Proof-of-Work in Bitum.
var MagicBytes = []byte{
	0x3b, 0xe5, 0xd4, 0x9e, 0xca, 0x59, 0x81, 0x5b,
	0x7a, 0x5d, 0xa1, 0xbb, 0x65, 0x37, 0x6f, 0x5d,
	0x04, 0xf4, 0x3b, 0x90, 0x2a, 0x41, 0x3a, 0xe5,
	0x8f, 0x98, 0x87, 0x93, 0xf9, 0x41, 0x67, 0x35,
}
