// Copyright (c) 2019 The Bitum developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

// bitumchain wraps codechain to handle the mainnet and testnet repositories.
package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/frankbraun/codechain/hashchain"
)

const (
	mainnetDir = ".codechain_mainnet"
	testnetDir = ".codechain_testnet"
)

var treehash = "undefined"

func callCodechain(args []string) error {
	cmd := exec.Command("codechain", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func lastTreeHash(dir string) (string, error) {
	c, err := hashchain.ReadFile(filepath.Join(dir, "hashchain"))
	if err != nil {
		return "", err
	}
	defer c.Close()
	return c.LastTreeHash(), nil
}

func readSourceLineComment() (string, error) {
	buf, err := ioutil.ReadFile(filepath.Join(testnetDir, "hashchain"))
	if err != nil {
		return "", err
	}
	lines := bytes.Split(buf, []byte("\n"))
	line := bytes.SplitN(lines[len(lines)-2], []byte(" "), 7)
	return string(line[6]), nil
}

func run(mainnet, testnet bool, args []string) error {
	fmt.Printf("_binary_treehash=%s\n", treehash)

	// check if we are publishing on both nets
	var sameOrigin bool
	if !mainnet && !testnet && args[0] == "publish" {
		testnetTreeHash, err := lastTreeHash(testnetDir)
		if err != nil {
			return err
		}
		fmt.Printf("testnet_treehash=%s\n", testnetTreeHash)
		mainnetTreeHash, err := lastTreeHash(mainnetDir)
		if err != nil {
			return err
		}
		fmt.Printf("mainnet_treehash=%s\n", mainnetTreeHash)
		if testnetTreeHash == mainnetTreeHash {
			sameOrigin = true
		}
	}

	// run codechain for testnet first
	if testnet || !mainnet {
		fmt.Printf("--------------------------------------------------------------------------------\n")
		fmt.Printf("$ env CODECHAIN_DIR=%s CODECHAIN_EXCLUDE=%s codechain %s\n",
			testnetDir, mainnetDir, strings.Join(args, " "))
		if err := os.Setenv("CODECHAIN_DIR", testnetDir); err != nil {
			return err
		}
		if err := os.Setenv("CODECHAIN_EXCLUDE", mainnetDir); err != nil {
			return err
		}
		if err := callCodechain(args); err != nil {
			return err
		}
	}
	// now run codechain for mainnet
	if mainnet || !testnet {
		if args[0] == "publish" {
			if sameOrigin {
				// we just published a new source version and both previous
				// versions are the same, read comment from hashchain
				msg, err := readSourceLineComment()
				if err != nil {
					return err
				}
				// publish with same comment on mainnet
				args = append(args, "-y", "-m", msg)
			} else {
				fmt.Println("originating tree hashes differ, review again")
			}
		}
		fmt.Printf("--------------------------------------------------------------------------------\n")
		fmt.Printf("$ env CODECHAIN_DIR=%s CODECHAIN_EXCLUDE=%s codechain %s\n",
			mainnetDir, testnetDir, strings.Join(args, " "))
		if err := os.Setenv("CODECHAIN_DIR", mainnetDir); err != nil {
			return err
		}
		if err := os.Setenv("CODECHAIN_EXCLUDE", testnetDir); err != nil {
			return err
		}
		if err := callCodechain(args); err != nil {
			return err
		}
	}
	return nil
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "%s: error: %s\n", os.Args[0], err)
	os.Exit(1)
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [options] codechain_command [codechain_options]\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	mainnet := flag.Bool("mainnet", false, "Call codechain only for mainnet")
	testnet := flag.Bool("testnet", false, "Call codechain only for testnet")
	flag.Usage = usage
	flag.Parse()
	if flag.NArg() < 1 {
		usage()
	}
	if *mainnet && *testnet {
		fatal(errors.New("-mainnet and -testnet exclude each other"))
	}
	if err := run(*mainnet, *testnet, flag.Args()); err != nil {
		fatal(err)
	}
}
