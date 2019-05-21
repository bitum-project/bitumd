// Copyright (c) 2019 The Bitum developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"encoding/hex"
	"testing"
)

var expectedScript []byte

func init() {
	var err error
	expectedScript, err = hex.DecodeString("a914cbb08d6ca783b533b2c7d24a51fbca92d937bf9987")
	if err != nil {
		panic(err)
	}
}

func TestAddr2PKScript(t *testing.T) {
	pkScript, err := addr2PKScript("ScuQxvveKGfpG1ypt6u27F99Anf7EW3cqhq")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(pkScript, expectedScript) {
		t.Errorf("Failed to convert address to public key script; "+
			"want %x, got %x", expectedScript, pkScript)
	}
}
