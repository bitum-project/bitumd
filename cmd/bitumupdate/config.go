// Copyright (c) 2018 The Decred developers
// Copyright (c) 2019 The Bitum developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	flags "github.com/jessevdk/go-flags"
	"github.com/bitum-project/bitumd/chaincfg"
	"github.com/bitum-project/bitumd/bitumutil"
)

const (
	defaultConfigFilename = "bitumupdate.conf"
)

var (
	// Default network parameters
	activeNetParams = &chaincfg.MainNetParams

	// Default configuration options
	defaultConfigFile = filepath.Join(defaultHomeDir, defaultConfigFilename)
	defaultHomeDir    = bitumutil.AppDataDir("bitumupdate", false)
)

// config defines the configuration options for bitumupdate.
//
// See loadConfig for details on the configuration load process.
type config struct {
	Node    string `short:"n" long:"node" description:"IP address of a  working node"`
	TestNet bool   `long:"testnet" description:"Use the test network"`
}

func loadConfig() (*config, error) {
	err := os.MkdirAll(defaultHomeDir, 0700)
	if err != nil {
		// Show a nicer error message if it's because a symlink is
		// linked to a directory that does not exist (probably because
		// it's not mounted).
		if e, ok := err.(*os.PathError); ok && os.IsExist(err) {
			if link, lerr := os.Readlink(e.Path); lerr == nil {
				str := "is symlink %s -> %s mounted?"
				err = fmt.Errorf(str, e.Path, link)
			}
		}

		str := "failed to create home directory: %v"
		err := fmt.Errorf(str, err)
		return nil, err
	}

	// Default config.
	cfg := config{}

	preCfg := cfg
	preParser := flags.NewParser(&preCfg, flags.Default)
	_, err = preParser.Parse()
	if err != nil {
		e, ok := err.(*flags.Error)
		if ok && e.Type == flags.ErrHelp {
			os.Exit(0)
		}
		preParser.WriteHelp(os.Stderr)
		return nil, err
	}

	appName := filepath.Base(os.Args[0])
	appName = strings.TrimSuffix(appName, filepath.Ext(appName))
	usageMessage := fmt.Sprintf("use %s -h to show usage", appName)

	// Load additional config from file.
	parser := flags.NewParser(&cfg, flags.Default)
	err = flags.NewIniParser(parser).ParseFile(defaultConfigFile)
	if err != nil {
		if _, ok := err.(*os.PathError); !ok {
			fmt.Fprintf(os.Stderr, "error parsing config "+
				"file: %v\n", err)
			fmt.Fprintln(os.Stderr, usageMessage)
			return nil, err
		}
	}

	// Parse command line options again to ensure they take precedence.
	_, err = parser.Parse()
	if err != nil {
		if e, ok := err.(*flags.Error); !ok || e.Type != flags.ErrHelp {
			parser.WriteHelp(os.Stderr)
		}
		return nil, err
	}

	if len(cfg.Node) == 0 {
		return nil, errors.New("please specify a node")
	}

	if cfg.TestNet {
		activeNetParams = &chaincfg.TestNetParams
	}

	return &cfg, nil
}
