// Copyright (c) 2019 The Bitum developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package main

import (
	"crypto/rand"
	"encoding/binary"
	"flag"
	"fmt"
	"os"
)

func gen32BitNonce() {
	b := make([]byte, 4)
	if _, err := rand.Read(b); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	nonce := binary.BigEndian.Uint32(b[0:])
	fmt.Printf("%#08x\n", nonce)
}

func gen256BitNonce() {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("{")
	for i, b := range buf {
		if i > 0 {
			fmt.Printf(", ")
		}
		fmt.Printf("%#02x", b)
	}
	fmt.Printf("}\n")
}

func main() {
	bit256 := flag.Bool("256", false, "generate 256-bit nonce instead of 32-bit")
	flag.Parse()
	if *bit256 {
		gen256BitNonce()
	} else {
		gen32BitNonce()
	}
}
