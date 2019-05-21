### Table of Contents
1. [About](#About)
2. [Getting Started](#GettingStarted)
    1. [Installation](#Installation)
    2. [Configuration](#Configuration)
    3. [Controlling and Querying bitumd via bitumctl](#DcrctlConfig)
    4. [Mining](#Mining)
3. [Help](#Help)
    1. [Network Configuration](#NetworkConfig)
    2. [Wallet](#Wallet)
4. [Contact](#Contact)
    1. [Community](#ContactCommunity)
5. [Developer Resources](#DeveloperResources)
    1. [Code Contribution Guidelines](#ContributionGuidelines)
    2. [JSON-RPC Reference](#JSONRPCReference)
    3. [Go Modules](#GoModules)
    4. [Module Hierarchy](#ModuleHierarchy)

<a name="About" />

### 1. About

bitumd is a full node Bitum implementation written in [Go](http://golang.org),
and is licensed under the [copyfree](http://www.copyfree.org) ISC License.

This software is currently under active development.  It is extremely stable and
has been in production use since February 2016.

It also properly relays newly mined blocks, maintains a transaction pool, and
relays individual transactions that have not yet made it into a block.  It
ensures all individual transactions admitted to the pool follow the rules
required into the block chain and also includes the vast majority of the more
strict checks which filter transactions based on miner requirements ("standard"
transactions).

<a name="GettingStarted" />

### 2. Getting Started

<a name="Installation" />

**2.1 Installation**<br />

The first step is to install bitumd.  The installation instructions can be found
[here](https://github.com/bitum-project/bitumd/tree/master/README.md#Installation).

<a name="Configuration" />

**2.2 Configuration**<br />

bitumd has a number of [configuration](http://godoc.org/github.com/bitum-project/bitumd)
options, which can be viewed by running: `$ bitumd --help`.

<a name="DcrctlConfig" />

**2.3 Controlling and Querying bitumd via bitumctl**<br />

bitumctl is a command line utility that can be used to both control and query bitumd
via [RPC](http://www.wikipedia.org/wiki/Remote_procedure_call).  bitumd does
**not** enable its RPC server by default;  You must configure at minimum both an
RPC username and password or both an RPC limited username and password:

* bitumd.conf configuration file
```
[Application Options]
rpcuser=myuser
rpcpass=SomeDecentp4ssw0rd
rpclimituser=mylimituser
rpclimitpass=Limitedp4ssw0rd
```
* bitumctl.conf configuration file
```
[Application Options]
rpcuser=myuser
rpcpass=SomeDecentp4ssw0rd
```
OR
```
[Application Options]
rpclimituser=mylimituser
rpclimitpass=Limitedp4ssw0rd
```
For a list of available options, run: `$ bitumctl --help`

<a name="Mining" />

**2.4 Mining**<br />
bitumd supports the [getwork](https://github.com/bitum-project/bitumd/tree/master/docs/json_rpc_api.md#getwork)
RPC.  The limited user cannot access this RPC.<br />

**1. Add the payment addresses with the `miningaddr` option.**<br />

```
[Application Options]
rpcuser=myuser
rpcpass=SomeDecentp4ssw0rd
miningaddr=DsExampleAddress1
miningaddr=DsExampleAddress2
```

**2. Add bitumd's RPC TLS certificate to system Certificate Authority list.**<br />

`cgminer` uses [curl](http://curl.haxx.se/) to fetch data from the RPC server.
Since curl validates the certificate by default, we must install the `bitumd` RPC
certificate into the default system Certificate Authority list.

**Ubuntu**<br />

1. Copy rpc.cert to /usr/share/ca-certificates: `# cp /home/user/.bitumd/rpc.cert /usr/share/ca-certificates/bitumd.crt`<br />
2. Add bitumd.crt to /etc/ca-certificates.conf: `# echo bitumd.crt >> /etc/ca-certificates.conf`<br />
3. Update the CA certificate list: `# update-ca-certificates`<br />

**3. Set your mining software url to use https.**<br />

`$ cgminer -o https://127.0.0.1:9209 -u rpcuser -p rpcpassword`

<a name="Help" />

### 3. Help

<a name="NetworkConfig" />

**3.1 Network Configuration**<br />
* [What Ports Are Used by Default?](https://github.com/bitum-project/bitumd/tree/master/docs/default_ports.md)
* [How To Listen on Specific Interfaces](https://github.com/bitum-project/bitumd/tree/master/docs/configure_peer_server_listen_interfaces.md)
* [How To Configure RPC Server to Listen on Specific Interfaces](https://github.com/bitum-project/bitumd/tree/master/docs/configure_rpc_server_listen_interfaces.md)
* [Configuring bitumd with Tor](https://github.com/bitum-project/bitumd/tree/master/docs/configuring_tor.md)

<a name="Wallet" />

**3.2 Wallet**<br />

bitumd was intentionally developed without an integrated wallet for security
reasons.  Please see [bitumwallet](https://github.com/bitum-project/bitumwallet) for more
information.

<a name="Contact" />

### 4. Contact

<a name="ContactCommunity" />

**4.1 Community**<br />

If you have any further questions you can find us at:

https://bitum.io/community

<a name="DeveloperResources" />

### 5. Developer Resources

<a name="ContributionGuidelines" />

**5.1 Code Contribution Guidelines**

* [Code Contribution Guidelines](https://github.com/bitum-project/bitumd/tree/master/docs/code_contribution_guidelines.md)

<a name="JSONRPCReference" />

**5.2 JSON-RPC Reference**

* [JSON-RPC Reference](https://github.com/bitum-project/bitumd/tree/master/docs/json_rpc_api.md)
    * [RPC Examples](https://github.com/bitum-project/bitumd/tree/master/docs/json_rpc_api.md#ExampleCode)

<a name="GoModules" />

**5.3 Go Modules**

The following versioned modules are provided by bitumd repository:

* [rpcclient](https://github.com/bitum-project/bitumd/tree/master/rpcclient) - Implements
  a robust and easy to use Websocket-enabled Bitum JSON-RPC client
* [bitumjson](https://github.com/bitum-project/bitumd/tree/master/bitumjson) - Provides an
  extensive API for the underlying JSON-RPC command and return values
* [wire](https://github.com/bitum-project/bitumd/tree/master/wire) - Implements the
  Bitum wire protocol
* [peer](https://github.com/bitum-project/bitumd/tree/master/peer) - Provides a common
  base for creating and managing Bitum network peers
* [blockchain](https://github.com/bitum-project/bitumd/tree/master/blockchain) -
  Implements Bitum block handling and chain selection rules
  * [stake](https://github.com/bitum-project/bitumd/tree/master/blockchain/stake) -
    Provides an API for working with stake transactions and other portions
    related to the Proof-of-Stake (PoS) system
* [txscript](https://github.com/bitum-project/bitumd/tree/master/txscript) -
  Implements the Bitum transaction scripting language
* [bitumec](https://github.com/bitum-project/bitumd/tree/master/bitumec) - Provides constants
  for the supported cryptographic signatures supported by Bitum scripts
  * [secp256k1](https://github.com/bitum-project/bitumd/tree/master/bitumec/secp256k1) -
    Implements the secp256k1 elliptic curve
  * [edwards](https://github.com/bitum-project/bitumd/tree/master/bitumec/edwards) -
    Implements the edwards25519 twisted Edwards curve
* [database](https://github.com/bitum-project/bitumd/tree/master/database) -
  Provides a database interface for the Bitum block chain
* [mempool](https://github.com/bitum-project/bitumd/tree/master/mempool) - Provides a
  policy-enforced pool of unmined Bitum transactions
* [bitumutil](https://github.com/bitum-project/bitumd/tree/master/bitumutil) - Provides
  Bitum-specific convenience functions and types
* [chaincfg](https://github.com/bitum-project/bitumd/tree/master/chaincfg) - Defines
  chain configuration parameters for the standard Bitum networks and allows
  callers to define their own custom Bitum networks for testing puproses
  * [chainhash](https://github.com/bitum-project/bitumd/tree/master/chaincfg/chainhash) -
    Provides a generic hash type and associated functions that allows the
    specific hash algorithm to be abstracted
* [certgen](https://github.com/bitum-project/bitumd/tree/master/certgen) - Provides a
  function for creating a new TLS certificate key pair, typically used for
  encrypting RPC and websocket communications
* [addrmgr](https://github.com/bitum-project/bitumd/tree/master/addrmgr) - Provides a
  concurrency safe Bitum network address manager
* [connmgr](https://github.com/bitum-project/bitumd/tree/master/connmgr) - Implements a
  generic Bitum network connection manager
* [hdkeychain](https://github.com/bitum-project/bitumd/tree/master/hdkeychain) - Provides
  an API for working with  Bitum hierarchical deterministic extended keys
* [gcs](https://github.com/bitum-project/bitumd/tree/master/gcs) - Provides an API for
  building and using Golomb-coded set filters useful for light clients such as
  SPV wallets
* [fees](https://github.com/bitum-project/bitumd/tree/master/fees) - Provides methods for
  tracking and estimating fee rates for new transactions to be mined into the
  network

<a name="ModuleHierarchy" />

**5.4 Module Hierarchy**

The following diagram shows an overview of the hierarchy for the modules
provided by the bitumd repository.

![Module Hierarchy](./assets/module_hierarchy.svg)
