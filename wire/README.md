wire
====

[![Build Status](http://img.shields.io/travis/bitum/bitumd.svg)](https://travis-ci.org/bitum/bitumd)
[![ISC License](http://img.shields.io/badge/license-ISC-blue.svg)](http://copyfree.org)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg)](http://godoc.org/github.com/bitum-project/bitumd/wire)

Package wire implements the Bitum wire protocol.  A comprehensive suite of
tests with 100% test coverage is provided to ensure proper functionality.

This package has intentionally been designed so it can be used as a standalone
package for any projects needing to interface with Bitum peers at the wire
protocol level.

## Installation and Updating

```bash
$ go get -u github.com/bitum-project/bitumd/wire
```

## Bitum Message Overview

The Bitum protocol consists of exchanging messages between peers. Each message
is preceded by a header which identifies information about it such as which
bitum network it is a part of, its type, how big it is, and a checksum to
verify validity. All encoding and decoding of message headers is handled by this
package.

To accomplish this, there is a generic interface for Bitum messages named
`Message` which allows messages of any type to be read, written, or passed
around through channels, functions, etc. In addition, concrete implementations
of most of the currently supported Bitum messages are provided. For these
supported messages, all of the details of marshalling and unmarshalling to and
from the wire using Bitum encoding are handled so the caller doesn't have to
concern themselves with the specifics.

## Reading Messages Example

In order to unmarshal Bitum messages from the wire, use the `ReadMessage`
function. It accepts any `io.Reader`, but typically this will be a `net.Conn`
to a remote node running a Bitum peer.  Example syntax is:

```Go
	// Use the most recent protocol version supported by the package and the
	// main Bitum network.
	pver := wire.ProtocolVersion
	bitumnet := wire.MainNet

	// Reads and validates the next Bitum message from conn using the
	// protocol version pver and the Bitum network bitumnet.  The returns
	// are a wire.Message, a []byte which contains the unmarshalled
	// raw payload, and a possible error.
	msg, rawPayload, err := wire.ReadMessage(conn, pver, bitumnet)
	if err != nil {
		// Log and handle the error
	}
```

See the package documentation for details on determining the message type.

## Writing Messages Example

In order to marshal Bitum messages to the wire, use the `WriteMessage`
function. It accepts any `io.Writer`, but typically this will be a `net.Conn`
to a remote node running a Bitum peer. Example syntax to request addresses
from a remote peer is:

```Go
	// Use the most recent protocol version supported by the package and the
	// main Bitum network.
	pver := wire.ProtocolVersion
	bitumnet := wire.MainNet

	// Create a new getaddr Bitum message.
	msg := wire.NewMsgGetAddr()

	// Writes a Bitum message msg to conn using the protocol version
	// pver, and the Bitum network bitumnet.  The return is a possible
	// error.
	err := wire.WriteMessage(conn, msg, pver, bitumnet)
	if err != nil {
		// Log and handle the error
	}
```

## License

Package wire is licensed under the [copyfree](http://copyfree.org) ISC
License.
