// Copyright (c) 2013-2016 The btcsuite developers
// Copyright (c) 2015-2016 The Decred developers
// Copyright (c) 2019 The Bitum developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package wire

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/bitum-project/bitumd/chaincfg/chainhash"
)

// TestInvVectStringer tests the stringized output for inventory vector types.
func TestInvTypeStringer(t *testing.T) {
	tests := []struct {
		in   InvType
		want string
	}{
		{InvTypeError, "ERROR"},
		{InvTypeTx, "MSG_TX"},
		{InvTypeBlock, "MSG_BLOCK"},
		{InvTypeFilteredBlock, "MSG_FILTERED_BLOCK"},
		{InvTypeCodechainEntry, "MSG_CODECHAIN_ENTRY"},
		{InvTypePatch, "MSG_PATCH"},
		{0xffffffff, "Unknown InvType (4294967295)"},
	}

	t.Logf("Running %d tests", len(tests))
	for i, test := range tests {
		result := test.in.String()
		if result != test.want {
			t.Errorf("String #%d\n got: %s want: %s", i, result,
				test.want)
			continue
		}
	}

}

// TestInvVect tests the InvVect API.
func TestInvVect(t *testing.T) {
	ivType := InvTypeBlock
	hash := chainhash.Hash{}

	// Ensure we get the same payload and signature back out.
	iv := NewInvVect(ivType, &hash)
	if iv.Type != ivType {
		t.Errorf("NewInvVect: wrong type - got %v, want %v",
			iv.Type, ivType)
	}
	if !iv.Hash.IsEqual(&hash) {
		t.Errorf("NewInvVect: wrong hash - got %v, want %v",
			spew.Sdump(iv.Hash), spew.Sdump(hash))
	}

}

// TestInvVectWire tests the InvVect wire encode and decode for various
// protocol versions and supported inventory vector types.
func TestInvVectWire(t *testing.T) {
	// Block 203707 hash.
	hashStr := "3264bc2ac36a60840790ba1d475d01367e7c723da941069e9dc"
	baseHash, err := chainhash.NewHashFromStr(hashStr)
	if err != nil {
		t.Errorf("NewHashFromStr: %v", err)
	}

	// Codechain head
	hashStr = "b3be280acb4a5aa5cc8b1b0cd58900f88f5b6fbdc2017d3ca8efd8a26dd3fe59"
	codechainEntry, err := chainhash.NewHashFromStr(hashStr)
	if err != nil {
		t.Errorf("NewHashFromStr: %v", err)
	}

	// Codechain patch
	hashStr = "4b607a05e42694f851de25a480f5bc7a8d856851c0954b4a3866bfaa396f839a"
	patch, err := chainhash.NewHashFromStr(hashStr)
	if err != nil {
		t.Errorf("NewHashFromStr: %v", err)
	}

	// errInvVect is an inventory vector with an error.
	errInvVect := InvVect{
		Type: InvTypeError,
		Hash: chainhash.Hash{},
	}

	// errInvVectEncoded is the wire encoded bytes of errInvVect.
	errInvVectEncoded := []byte{
		0x00, 0x00, 0x00, 0x00, // InvTypeError
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // No hash
	}

	// txInvVect is an inventory vector representing a transaction.
	txInvVect := InvVect{
		Type: InvTypeTx,
		Hash: *baseHash,
	}

	// txInvVectEncoded is the wire encoded bytes of txInvVect.
	txInvVectEncoded := []byte{
		0x01, 0x00, 0x00, 0x00, // InvTypeTx
		0xdc, 0xe9, 0x69, 0x10, 0x94, 0xda, 0x23, 0xc7,
		0xe7, 0x67, 0x13, 0xd0, 0x75, 0xd4, 0xa1, 0x0b,
		0x79, 0x40, 0x08, 0xa6, 0x36, 0xac, 0xc2, 0x4b,
		0x26, 0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // Block 203707 hash
	}

	// blockInvVect is an inventory vector representing a block.
	blockInvVect := InvVect{
		Type: InvTypeBlock,
		Hash: *baseHash,
	}

	// blockInvVectEncoded is the wire encoded bytes of blockInvVect.
	blockInvVectEncoded := []byte{
		0x02, 0x00, 0x00, 0x00, // InvTypeBlock
		0xdc, 0xe9, 0x69, 0x10, 0x94, 0xda, 0x23, 0xc7,
		0xe7, 0x67, 0x13, 0xd0, 0x75, 0xd4, 0xa1, 0x0b,
		0x79, 0x40, 0x08, 0xa6, 0x36, 0xac, 0xc2, 0x4b,
		0x26, 0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // Block 203707 hash
	}

	// codechainEntryInvVect is an inventory vector representing a Codechain entry.
	codechainEntryInvVect := InvVect{
		Type: InvTypeCodechainEntry,
		Hash: *codechainEntry,
	}

	// codechainEntryInvVectEncoded is the wire encoded bytes of codechainEntryInvVect.
	codechainEntryInvVectEncoded := []byte{
		0x00, 0x01, 0x00, 0x00, // InvTypeCodechainEntry
		0x59, 0xfe, 0xd3, 0x6d, 0xa2, 0xd8, 0xef, 0xa8,
		0x3c, 0x7d, 0x01, 0xc2, 0xbd, 0x6f, 0x5b, 0x8f,
		0xf8, 0x00, 0x89, 0xd5, 0x0c, 0x1b, 0x8b, 0xcc,
		0xa5, 0x5a, 0x4a, 0xcb, 0x0a, 0x28, 0xbe, 0xb3, // Codechain head
	}

	// patchInvVect is an inventory vector representing a patch.
	patchInvVect := InvVect{
		Type: InvTypePatch,
		Hash: *patch,
	}

	// patchInvVectEncoded is the wire encoded bytes of patchInvVect.
	patchInvVectEncoded := []byte{
		0x01, 0x01, 0x00, 0x00, // InvTypePatch
		0x9a, 0x83, 0x6f, 0x39, 0xaa, 0xbf, 0x66, 0x38,
		0x4a, 0x4b, 0x95, 0xc0, 0x51, 0x68, 0x85, 0x8d,
		0x7a, 0xbc, 0xf5, 0x80, 0xa4, 0x25, 0xde, 0x51,
		0xf8, 0x94, 0x26, 0xe4, 0x05, 0x7a, 0x60, 0x4b, // Codechain patch
	}

	tests := []struct {
		in   InvVect // NetAddress to encode
		out  InvVect // Expected decoded NetAddress
		buf  []byte  // Wire encoding
		pver uint32  // Protocol version for wire encoding
	}{
		// Latest protocol version error inventory vector.
		{
			errInvVect,
			errInvVect,
			errInvVectEncoded,
			ProtocolVersion,
		},

		// Latest protocol version tx inventory vector.
		{
			txInvVect,
			txInvVect,
			txInvVectEncoded,
			ProtocolVersion,
		},

		// Latest protocol version block inventory vector.
		{
			blockInvVect,
			blockInvVect,
			blockInvVectEncoded,
			ProtocolVersion,
		},

		// Latest protocol version Codechain entry inventory vector.
		{
			codechainEntryInvVect,
			codechainEntryInvVect,
			codechainEntryInvVectEncoded,
			ProtocolVersion,
		},

		// Latest protocol version patch inventory vector.
		{
			patchInvVect,
			patchInvVect,
			patchInvVectEncoded,
			ProtocolVersion,
		},
	}

	t.Logf("Running %d tests", len(tests))
	for i, test := range tests {
		// Encode to wire format.
		var buf bytes.Buffer
		err := writeInvVect(&buf, test.pver, &test.in)
		if err != nil {
			t.Errorf("writeInvVect #%d error %v", i, err)
			continue
		}
		if !bytes.Equal(buf.Bytes(), test.buf) {
			t.Errorf("writeInvVect #%d\n got: %s want: %s", i,
				spew.Sdump(test.buf), spew.Sdump(buf.Bytes()))
			continue
		}

		// Decode the message from wire format.
		var iv InvVect
		rbuf := bytes.NewReader(test.buf)
		err = readInvVect(rbuf, test.pver, &iv)
		if err != nil {
			t.Errorf("readInvVect #%d error %v", i, err)
			continue
		}
		if !reflect.DeepEqual(iv, test.out) {
			t.Errorf("readInvVect #%d\n got: %s want: %s", i,
				spew.Sdump(test.out), spew.Sdump(iv))
			continue
		}
	}
}
