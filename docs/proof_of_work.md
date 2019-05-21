### Proof-of-Work

This document describes the Proof-of-Work (PoW) mechanism used in
Bitum.

The PoW in mechanism in Bitum is similar to the one employed in
Bitcoin, namely a double SHA-256. But additonally, after each round of
SHA-256 the resulting hash is XORed with a 256-bit magic number (defined
in `chainhash.MagicBytes`):

```
    sha256 (sha256 ( header (otherdata, nonce) ) ^ magic_number ) ^ magic_number < difficulty_threshold
```

This makes existing Bitcoin ASICs useless, but it would only take minor
changes to the hardware design to create Bitum ASICs.

For all other hash function calls in Bitum BLAKE-256 is retained.

See:

*   `$GOPATH/src/github.com/bitum-project/bitumd/cpuminer.go`
*   `$GOPATH/src/github.com/btcsuite/btcd/mining/cpuminer/cpuminer.go`
*   `$GOPATH/src/github.com/bitum-project/bitumd/wire/blockheader.go`
    (`BlockHash()`)
*   `$GOPATH/src/github.com/btcsuite/btcd/wire/blockheader.go`
    (`BlockHash()`)
*   `$GOPATH/src/github.com/bitum-project/bitumd/chaincfg/chainhash/hash.go`
*   `$GOPATH/src/github.com/btcsuite/btcd/chaincfg/chainhash/hash.go`


# Bitcoin Proof-of-Work

Roughly described, the proof-of-work protocol in bitcoin works as follows:

The bitcoin block header contains various pieces of information:

* block version
* previous block header hash
* root hash of merkle tree containing transactions
* time
* difficulty threshold (encoded as nbits)
* nonce


Proof of work functions by varying the nonce value until a valid block is found, which is defined as:

* Merkle root is a valid root for a merkle tree which itself only references transactions which are also valid.
* Block weight does not exceed the maximum.
* Difficulty threshold is correct given the difficulty adjustment algorithm.
* Previous block header hash references a previously valid block.
* Block version is higher or equal to previous block version.
* The hash of the block is valid according to the proof-of-work algorithm, ie:

```
    sha256 (sha256 ( header (otherdata, nonce) ) ) < difficulty threshold
```

# Bitum Proof-of-Work

The only change in Bitum would be the addition of two XOR operations in the proof-of-work algorithm. All valid blocks therefore must meet the following constraint:

```
    sha256 (sha256 ( header (otherdata, nonce) ) ^ magic_number ) ^ magic_number < difficulty_threshold
```

