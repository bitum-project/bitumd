### Bitum updater

This document describes the mechanism to publish Bitum updates.

We use [Codechain](https://github.com/frankbraun/codechain) to secure
our code against unwanted changes. The last signed source tree hash of
the Codechain hashchain is published in the Bitum block headers.

### Mechanism

-   Use two different Codechain hashchains, one for the Bitum
    mainnet and one for testnet (`.codechain_mainnet` and
    `.codechain_testnet`).
-   The mainnet and testnet Codechains have different sets of signers.
    The security requirements for testnet are less strict.
-   `bitumchain` is a Codechain wrapper that calls `codechain` for
    mainnet and testnet. `codechain publish -y` is used internally to
    avoid reading the same patch twice.
-   Extend `wire.BlockHeader` with a `CodechainHead [32]byte` field to
    publish the current Codechain source tree hash used by the miner.
    Maybe `ExtraData [32]byte` could also be used for that.
-   Extend the `wire` format to distribute Codechain hashchain updates
    and patch files (details TBD):
    -   Give me hashchain updates (between `X` and `Y`)
    -   Give me entire hashchain
    -   Give me patch
    -   Similar methods for publishing (via `bitumupdate`)?
-   Miners manually update to a newly published version. As soon as a
    miner mines a block with the new header it _activates_ on the
    network. Miners always write the version they are running on into
    the header.
-   Nodes either update automatically or manually. Automatic updates
    happen probabilistically (10% probability that a node updates in a
    given day, with manual override). This leads to a good code
    diversity in the network even in the case of automatic updates.
-   It is important that not every node updates to the newest version
    immediately, there should always be `N` versions of the code running
    on the network (with `N` between 3 and 5).
-   `bitumd` has the corresponding `Codechain` hash compiled in and it
    verifies that it is a valid one and in the allowed window defined by
    `N`.
-   `bitumupdate` is used to publish an update to the network via the
    extended the `wire` protocol.

### Wire protocol

-   `MsgVersion` & `MsgVerAck`
-   Encode current code version in `MsgVersion.UserAgent` or in separate
    field?
