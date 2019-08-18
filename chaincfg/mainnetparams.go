// Copyright (c) 2014-2016 The btcsuite developers
// Copyright (c) 2015-2019 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package chaincfg

import (
	"time"

	"github.com/bitum-project/bitumd/wire"
)

// MainNetParams defines the network parameters for the main Bitum network.
var MainNetParams = Params{
	Name:        "mainnet",
	Net:         wire.MainNet,
	DefaultPort: "9208",
	DNSSeeds: []DNSSeed{
		{"dnsseed.bitum.io", true},
	},

	// Chain parameters
	GenesisBlock:             &genesisBlock,
	GenesisHash:              &genesisHash,
	PowLimit:                 mainPowLimit,
	PowLimitBits:             0x1d00ffff,
	ReduceMinDifficulty:      false,
	MinDiffReductionTime:     0, // ~99.3% chance to be mined before reduction
	GenerateSupported:        true,
	MaximumBlockSizes:        []int{393216},
	MaxTxSize:                393216,
	TargetTimePerBlock:       time.Minute * 5,
	WorkDiffAlpha:            1,
	WorkDiffWindowSize:       144,
	WorkDiffWindows:          20,
	TargetTimespan:           time.Minute * 5 * 144,
	RetargetAdjustmentFactor: 4,

	// Subsidy parameters.
	BaseSubsidy:              3119582664,
	MulSubsidy:               100,
	DivSubsidy:               101,
	SubsidyReductionInterval: 6144,
	WorkRewardProportion:     6,
	StakeRewardProportion:    3,
	BlockTaxProportion:       1,

	// Checkpoints ordered from oldest to newest.
	Checkpoints: []Checkpoint{
		{100, newHashFromStr("000000005b24df7dd3dcdfbb4a90e4001963360b4181f4975e9b94a3d94039a8")},
		{500, newHashFromStr("000000001d5f634c9fda95180ccb472de9cbc7d25e3fea276b8b2706ea04a610")},
		{1000, newHashFromStr("000000000037c4bba623aa717b50e530a5f9fd891df815e2791cd0a3a233b782")},
		{2000, newHashFromStr("0000000000104f477a38499a5988c5ace7e155e9fb27554b955f3e22724736cc")},
		{3000, newHashFromStr("00000000000207e68b97cf74585aec083d3118a524f50a177615622bf0bb2b9c")},
		{4000, newHashFromStr("00000000001461c0cc9e88eca5a2e82029dddf240457d3a4b99725984a8362c2")},
		{9000, newHashFromStr("000000000000af2d102346b800d7b9fb9c9cfa71677fd3bcd77eb7b03d20a290")},
		{9583, newHashFromStr("000000000030279de3cc16ac237f264471d44e75c89efae4e9add41e9c50c0a5")},
	},

	// The miner confirmation window is defined as:
	//   target proof of work timespan / target proof of work spacing
	RuleChangeActivationQuorum:     4032, // 10 % of RuleChangeActivationInterval * TicketsPerBlock
	RuleChangeActivationMultiplier: 3,    // 75%
	RuleChangeActivationDivisor:    4,
	RuleChangeActivationInterval:   2016 * 4, // 4 weeks
	
	Deployments: map[uint32][]ConsensusDeployment{},
	
	// Enforce current block version once majority of the network has
	// upgraded.
	// 75% (750 / 1000)
	// Reject previous block versions once a majority of the network has
	// upgraded.
	// 95% (950 / 1000)
	BlockEnforceNumRequired: 750,
	BlockRejectNumRequired:  950,
	BlockUpgradeNumToCheck:  1000,

	// AcceptNonStdTxs is a mempool param to either accept and relay
	// non standard txs to the network or reject them
	AcceptNonStdTxs: false,

	// Address encoding magics
	NetworkAddressPrefix: "B",
	PubKeyAddrID:         [2]byte{0x11, 0x86},
	PubKeyHashAddrID:     [2]byte{0x05, 0xa3},
	PKHEdwardsAddrID:     [2]byte{0x09, 0x1f},
	PKHSchnorrAddrID:     [2]byte{0x08, 0x01},
	ScriptHashAddrID:     [2]byte{0x07, 0x1a},
	PrivateKeyID:         [2]byte{0x06, 0xde},

	// BIP32 hierarchical deterministic extended key magics
	HDPrivateKeyID: [4]byte{0x01, 0xf3, 0xa5, 0xe3},
	HDPublicKeyID:  [4]byte{0x02, 0xf1, 0xa7, 0x17},

	// BIP44 coin type used in the hierarchical deterministic path for
	// address generation.
	SLIP0044CoinType: 42, // SLIP0044, Bitum
	LegacyCoinType:   20, // for backwards compatibility

	// Bitum PoS parameters
	MinimumStakeDiff:        2 * 1e8, // 2 Coin
	TicketPoolSize:          8192,
	TicketsPerBlock:         5,
	TicketMaturity:          256,
	TicketExpiry:            40960, // 5*TicketPoolSize
	CoinbaseMaturity:        256,
	SStxChangeMaturity:      1,
	TicketPoolSizeWeight:    4,
	StakeDiffAlpha:          1, // Minimal
	StakeDiffWindowSize:     144,
	StakeDiffWindows:        20,
	StakeVersionInterval:    144 * 2 * 7, // ~1 week
	MaxFreshStakePerBlock:   20,          // 4*TicketsPerBlock
	StakeEnabledHeight:      256 + 256,   // CoinbaseMaturity + TicketMaturity
	StakeValidationHeight:   4096,        // ~14 days
	StakeBaseSigScript:      []byte{0x00, 0x00},
	StakeMajorityMultiplier: 3,
	StakeMajorityDivisor:    4,

	// Bitum organization related parameters
	// Organization address is B1xAWYg2eAyXhbetkLTMWmWN3Ub8AZfkeTq
	OrganizationPkScript:        hexDecode("76a914ca62b11e8a5ca4ea64616604adf12c819cfcc3f788ac"),
	OrganizationPkScriptVersion: 0,
	BlockOneLedger:              BlockOneLedgerMainNet,
}
