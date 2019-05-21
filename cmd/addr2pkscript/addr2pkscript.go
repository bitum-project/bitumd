// Copyright (c) 2019 The Bitum developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

// addr2pkscript converts a Bitum address to a public key script.
package main

import (
	"fmt"
	"os"

	"github.com/bitum-project/bitumd/txscript"
	"github.com/bitum-project/bitumd/bitumutil"
)

func addr2PKScript(addrStr string) ([]byte, error) {
	addr, err := bitumutil.DecodeAddress(addrStr)
	if err != nil {
		return nil, err
	}
	pkScript, err := txscript.PayToAddrScript(addr)
	if err != nil {
		return nil, err
	}
	return pkScript, nil
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "%s: error: %s\n", os.Args[0], err)
	os.Exit(1)
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s address\n", os.Args[0])
	os.Exit(2)
}

func main() {
	if len(os.Args) != 2 {
		usage()
	}
	pkScript, err := addr2PKScript(os.Args[1])
	if err != nil {
		fatal(err)
	}
	fmt.Printf("%x\n", pkScript)
}
