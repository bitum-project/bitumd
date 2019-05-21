// Copyright (c) 2018 The Decred developers
// Copyright (c) 2019 The Bitum developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

// bitumupdate pushes signed Bitum updates to the network.
package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/bitum-project/bitumd/chaincfg/chainhash"
	"github.com/bitum-project/bitumd/peer"
	"github.com/bitum-project/bitumd/updater"
	"github.com/bitum-project/bitumd/wire"
)

const (
	// defaultNodeTimeout defines the timeout time waiting for
	// a response from a node.
	defaultNodeTimeout = time.Second * 10
)

func newMsgInvCodechainEntry(head [32]byte) *wire.MsgInv {
	var h chainhash.Hash
	copy(h[:], head[:])
	invMsg := wire.NewMsgInvSizeHint(1)
	iv := wire.NewInvVect(wire.InvTypeCodechainEntry, &h)
	err := invMsg.AddInvVect(iv)
	if err != nil {
		panic(err)
	}
	return invMsg
}

func pushUpdate(node string, head [32]byte) error {
	// Connect to peer node.
	verack := make(chan struct{})
	config := peer.Config{
		UserAgentName:  "bitumupdate",
		ChainParams:    activeNetParams,
		DisableRelayTx: true,
		Listeners: peer.MessageListeners{
			OnVerAck: func(p *peer.Peer, msg *wire.MsgVerAck) {
				log.Printf("adding peer %v with services %v",
					p.NA().IP.String(), p.Services())
				verack <- struct{}{}
			},
		},
	}
	host := net.JoinHostPort(node, activeNetParams.DefaultPort)
	p, err := peer.NewOutboundPeer(&config, host)
	if err != nil {
		return err
	}
	conn, err := net.DialTimeout("tcp", p.Addr(), defaultNodeTimeout)
	if err != nil {
		return err
	}
	p.AssociateConnection(conn)

	// Wait for the verack message or timeout in case of failure.
	select {
	case <-verack:
		log.Printf("verack on peer %v", p.Addr())
	case <-time.After(defaultNodeTimeout):
		err := fmt.Errorf("verack timeout on peer %v", p.Addr())
		p.Disconnect()
		return err
	}

	// Send current head as inventory message to peer.
	done := make(chan struct{})
	p.QueueMessage(newMsgInvCodechainEntry(head), done)
	<-done
	log.Print("inv message sent")

	// TODO

	p.Disconnect()
	return nil
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "%s: error: %s\n", os.Args[0], err)
	os.Exit(1)
}

func main() {
	cfg, err := loadConfig()
	if err != nil {
		fatal(err)
	}
	head, err := updater.Head(activeNetParams.Net)
	if err != nil {
		fatal(err)
	}
	log.Print("head:")
	log.Printf("%x", *head)
	if err := pushUpdate(cfg.Node, *head); err != nil {
		fatal(err)
	}
}
