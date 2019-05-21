// Copyright (c) 2013-2016 The btcsuite developers
// Copyright (c) 2015-2018 The Decred developers
// Copyright (c) 2019 The Bitum developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package updater

import (
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"

	"github.com/frankbraun/codechain/hashchain"
	"github.com/frankbraun/codechain/util/file"
	"github.com/frankbraun/codechain/util/hex"
	"github.com/bitum-project/bitumd/peer"
	"github.com/bitum-project/bitumd/updater/internal/chainstate"
	"github.com/bitum-project/bitumd/wire"
)

const (
	mainnetHeadStr = "4fc99b1f42c57c3e0918162006ec85b8fbe2e0e25cf9f36ad287a9f56d2463e9"
	testnetHeadStr = "246d7e6ab314d694fd602b04d5aa46bd6713e667e0e1ad6675533787b2b3fe00"
)

// invMsg packages a Bitum inv message and the peer it came from together
// so the update handler has access to that information.
type invMsg struct {
	inv  *wire.MsgInv
	peer *peer.Peer
}

// UpdateManager provides a concurrency safe update manager for handling all
// incoming Codechain entries and patch files.
type UpdateManager struct {
	started    int32
	shutdown   int32
	msgChan    chan interface{}
	wg         sync.WaitGroup
	quit       chan struct{}
	chainState *chainstate.ChainState
}

func containsStringAtPos(sa []string, s string) (bool, int) {
	for i, v := range sa {
		if v == s {
			return true, i
		}
	}
	return false, -1
}

func checkHashchain(filename, treeHash string) error {
	hc, err := hashchain.ReadFile(filename)
	if err != nil {
		return err
	}
	defer hc.Close()

	// Make sure hash chain contains known head.
	var (
		mainnetHead [32]byte
		testnetHead [32]byte
	)
	head, err := hex.Decode(mainnetHeadStr, 32)
	if err != nil {
		return err
	}
	copy(mainnetHead[:], head)
	head, err = hex.Decode(testnetHeadStr, 32)
	if err != nil {
		return err
	}
	copy(testnetHead[:], head)
	if err := hc.CheckHead(mainnetHead); err != nil {
		log.Warnf("Head for mainnet not found, checking testnet: %v", err)
		if err := hc.CheckHead(testnetHead); err != nil {
			return err
		}
		log.Info("Head for testnet found.")
	}

	// Try to find treeHash in hash chain (signed or unsigned).
	if treeHash == "undefined" {
		log.Warn("treehash is \"undefined\", abort check.")
		return nil
	}
	// Make sure the treehash is parsable.
	_, err = hex.Decode(treeHash, 32)
	if err != nil {
		return err
	}
	treeHashes := hc.TreeHashes()
	contains, pos := containsStringAtPos(treeHashes, treeHash)
	if !contains {
		log.Warnf("Treehash %s not found.", treeHash)
		return nil
	}
	_, lastPos := hc.LastSignedTreeHash()
	if pos <= lastPos {
		log.Infof("Signed treehash %s found.", treeHash)
	} else {
		log.Infof("Unsigned treehash %s found.", treeHash)
	}
	return nil
}

func checkTreeHash(treeHash, dataDir string) error {
	codechainDir := filepath.Join(dataDir, "src", ".codechain")
	log.Infof("Checking treehash %s in %s", treeHash, codechainDir)
	if err := os.MkdirAll(codechainDir, 0755); err != nil {
		return err
	}
	hashchainFile := filepath.Join(codechainDir, "hashchain")
	exists, err := file.Exists(hashchainFile)
	if err != nil {
		return err
	}
	if !exists {
		log.Warnf("%s: file does not exist, abort check.", hashchainFile)
		return nil
	}
	return checkHashchain(hashchainFile, treeHash)
}

// NewUpdateManager returns a new Bitum update manager.
// Use Start to begin processing asynchronous update inv updates.
func NewUpdateManager(maxPeers int, treeHash, dataDir string) (*UpdateManager, error) {
	um := UpdateManager{
		msgChan: make(chan interface{}, maxPeers*3),
		quit:    make(chan struct{}),
	}

	// Check source code version at startup.
	err := checkTreeHash(treeHash, dataDir)
	if err != nil {
		return nil, err
	}

	um.chainState, err = chainstate.New(dataDir)
	if err != nil {
		return nil, err
	}

	return &um, nil
}

// handleInvMsg handles inv messages from all peers.
// We examine the inventory advertised by the remote peer and act accordingly.
func (u *UpdateManager) handleInvMsg(imsg *invMsg) {
	log.Info("Handling inv message...")
	for _, invVect := range imsg.inv.InvList {
		switch invVect.Type {
		case wire.InvTypeCodechainEntry:
			log.Infof("Codechain entry: %x", invVect.Hash[:])
			if u.chainState.EntryIsKnown(hex.Encode(invVect.Hash[:])) {
				log.Info("Entry is known.")
			} else {
				log.Info("Entry is unknown.")
			}
		case wire.InvTypePatch:
			log.Infof("Patch file: %x", invVect.Hash[:])
		default:
			log.Warnf("Type not handled here: %s", invVect.Type)
		}
	}
}

// blockHandler is the main handler for the block manager.  It must be run
// as a goroutine.  It processes block and inv messages in a separate goroutine
// from the peer handlers so the block (MsgBlock) messages are handled by a
// single thread without needing to lock memory data structures.  This is
// important because the block manager controls which blocks are needed and how
// the fetching should proceed.
func (u *UpdateManager) updateHandler() {
out:
	for {
		select {
		case m := <-u.msgChan:
			switch msg := m.(type) {
			case *invMsg:
				u.handleInvMsg(msg)
			default:
				log.Warnf("Invalid message type in update "+
					"handler: %T", msg)
			}

		case <-u.quit:
			break out
		}
	}

	u.chainState.Close()
	u.wg.Done()
	log.Trace("Update handler done")
}

// QueueInv adds the passed inv message and peer to the update handling queue.
func (u *UpdateManager) QueueInv(inv *wire.MsgInv, p *peer.Peer) {
	// No channel handling here because peers do not need to block on inv
	// messages.
	if atomic.LoadInt32(&u.shutdown) != 0 {
		return
	}

	u.msgChan <- &invMsg{inv: inv, peer: p}
}

// Start begins the core update handler which processes Codechain entry and
// patch file inv messages.
func (u *UpdateManager) Start() {
	// Already started?
	if atomic.AddInt32(&u.started, 1) != 1 {
		return
	}

	log.Trace("Starting update manager")
	u.wg.Add(1)
	go u.updateHandler()
}

// Stop gracefully shuts down the update manager by stopping all asynchronous
// handlers and waiting for them to finish.
func (u *UpdateManager) Stop() error {
	if atomic.AddInt32(&u.shutdown, 1) != 1 {
		log.Warnf("Update manager is already in the process of " +
			"shutting down")
		return nil
	}

	log.Infof("Update manager shutting down")
	close(u.quit)
	u.wg.Wait()
	return nil
}
