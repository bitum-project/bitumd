fees
=======


[![Build Status](http://img.shields.io/travis/bitum/bitumd.svg)](https://travis-ci.org/bitum/bitumd)
[![ISC License](http://img.shields.io/badge/license-ISC-blue.svg)](http://copyfree.org)
[![GoDoc](http://img.shields.io/badge/godoc-reference-blue.svg)](http://godoc.org/github.com/bitum-project/bitumd/fees)

Package fees provides bitum-specific methods for tracking and estimating fee
rates for new transactions to be mined into the network. Fee rate estimation has
two main goals:

- Ensuring transactions are mined within a target _confirmation range_
  (expressed in blocks);
- Attempting to minimize fees while maintaining be above restriction.

This package was started in order to resolve issue bitum/bitumd#1412 and related.
See that issue for discussion of the selected approach.

This package was developed for bitumd, a full-node implementation of Bitum which
is under active development.  Although it was primarily written for
bitumd, this package has intentionally been designed so it can be used as a
standalone package for any projects needing the functionality provided.

## Installation and Updating

```bash
$ go get -u github.com/bitum-project/bitumd/fees
```

## License

Package bitumutil is licensed under the [copyfree](http://copyfree.org) ISC
License.
