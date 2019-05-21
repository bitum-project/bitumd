Bitum
=========

This repository contains `bitumd`, a Bitum full node implementation
written in Go (Golang).

It acts as a chain daemon for the [Bitum](https://bitum.io)
cryptocurrency. `bitumd` maintains the entire past transactional ledger
of Bitum and allows relaying of transactions to other Bitum
nodes across the world. To read more about Bitum please see the
[project documentation](/doc/trunk/docs/overview.md).

To send or receive funds and join Proof-of-Stake mining, you will also
need to use a `bitumwallet`.

Requirements
------------

[Go](http://golang.org) 1.11 or newer.

Installation
------------

```
go get -u -v github.com/bitum-project/bitumd/...
```


Testing
-------

`$ make test`

(Runs `./run_tests.sh`.)

Issue Tracker
-------------

The [integrated issue tracker](/ticket) is used for this project.

Documentation
-------------

The documentation is a work-in-progress. It is located in the
[docs](/dir?ci=trunk&name=docs) folder. See for example:

-   [Development notes](docs/development_notes.md)
-   [Bitum updater](docs/updater.md)
-   [Proof-of-Work mechanism](docs/proof_of_work.md)

License
-------

Bitum is licensed under the [copyfree](http://copyfree.org) ISC
License.
