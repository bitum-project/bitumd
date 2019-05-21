// Copyright (c) 2014 The btcsuite developers
// Copyright (c) 2015-2018 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package chaincfg

// BlockOneLedgerMainNet is the block one output ledger for the main
// network.
var BlockOneLedgerMainNet = []*TokenPayout{
//  {"FeRpwwMm5S5b6YF1RaLS3zhsZHdcr3QzbCg", 100000 * 1e8},
}

// BlockOneLedgerTestNet is the block one output ledger for the test
// network.
var BlockOneLedgerTestNet = []*TokenPayout{}

// BlockOneLedgerSimNet is the block one output ledger for the simulation
// network.  See "Bitum organization related parameters" in simnetparams.go for
// information on how to spend these outputs.
var BlockOneLedgerSimNet = []*TokenPayout{}

// BlockOneLedgerRegNet is the block one output ledger for the regression test
// network.  See "Bitum organization related parameters" in regnetparams.go for
// information on how to spend these outputs.
var BlockOneLedgerRegNet = []*TokenPayout{}
