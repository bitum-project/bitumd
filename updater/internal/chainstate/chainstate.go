// Copyright (c) 2019 The Bitum developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

// Package chainstate implements the hashchain state for the Bitum updater.
package chainstate

import (
	"bytes"
	"crypto/sha256"
	"path/filepath"

	"github.com/frankbraun/codechain/hashchain"
	"github.com/frankbraun/codechain/tree"
	"github.com/frankbraun/codechain/util/file"
	"github.com/frankbraun/codechain/util/hex"
)

// ChainState holds the hashchain state for the Bitum updater.
type ChainState struct {
	hashchainEntries map[string]bool
	head             [32]byte
	patchFiles       map[string]bool
}

// New returns a new ChainState.
func New(dataDir string) (*ChainState, error) {
	var cs ChainState
	cs.hashchainEntries = make(map[string]bool)
	cs.patchFiles = make(map[string]bool)

	hashchainFile := filepath.Join(dataDir, "src", ".codechain", "hashchain")
	exists, err := file.Exists(hashchainFile)
	if err != nil {
		return nil, err
	}

	// add empty hash
	log.Debug("Adding hashchain entries...")
	cs.hashchainEntries[tree.EmptyHash] = true
	head, err := hex.Decode(tree.EmptyHash, 32)
	if err != nil {
		return nil, err
	}
	copy(cs.head[:], head)
	log.Debug(tree.EmptyHash)
	if exists {
		// add other entries
		hc, err := hashchain.ReadFile(hashchainFile)
		if err != nil {
			return nil, err
		}
		defer hc.Close()
		var buf bytes.Buffer
		if err := hc.Fprint(&buf); err != nil {
			return nil, err
		}
		for _, line := range bytes.Split(buf.Bytes(), []byte("\n")) {
			hash := sha256.Sum256(line)
			hashStr := hex.Encode(hash[:])
			log.Debug(hashStr)
			cs.hashchainEntries[hashStr] = true
		}
		head := hc.Head()
		copy(cs.head[:], head[:])
	}

	return &cs, nil
}

// Close chain state.
func (cs *ChainState) Close() {
	// TODO
}

// EntryIsKnown returns truf the given hashchainEntry is known and false
// otherwise.
func (cs *ChainState) EntryIsKnown(hashchainEntry string) bool {
	return cs.hashchainEntries[hashchainEntry]
}

// Head returns the hash of the last known hashchain entry.
func (cs *ChainState) Head() [32]byte {
	return cs.head
}
