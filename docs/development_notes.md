Development notes
-----------------

### Runing a Simnet

See [Simnet](https://docs.decred.org/advanced/simnet/).

### Running a local Testnet

Daemon 1:

    bitumd --testnet --nodnsseed --generate \
           --miningaddr=Tsp7dmytXg96vk7MwYV9KSdf2UcWxjxRfPn

Daemon 2:

    bitumd --testnet --nodnsseed --generate \
           --nolisten --norpc --appdata=/home/user/.bitumd2 --connect=localhost \
           --miningaddr=Tsd92111UDM8H4pxTJb5LBgDJKxEzpgGiyu

Wallet:

    bitumwallet --testnet

The two nodes have to talk to be connected to each other before they
start mining.
