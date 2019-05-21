### Making transaction from multisig address

WARNING! In current documentation fee is not configured properly.

First, you have prepare variables for creating a raw transaction

```
utxo_txid=$(bitumctl --testnet --wallet listunspent | jq -r '.[0] | .txid')
utxo_vout=$(bitumctl --testnet --wallet listunspent | jq -r '.[0] | .vout')
utxo_tree=$(bitumctl --testnet --wallet listunspent | jq -r '.[0] | .vout')
utxo_spk=$(bitumctl --testnet --wallet listunspent | jq -r '.[0] | .scriptPubKey')

recipient="recipient_address"
redeem_script="redeem_script_code"
```

Than you create raw transaction by using this command
```
rawtxhex=$(bitumctl --testnet --wallet createrawtransaction "[{\"txid\": \"$utxo_txid\", \"vout\": $utxo_vout, \"tree\": $utxo_tree}]" "{\"$recipient\": 0.33}")
```

Next is it needed to get private key for first address, It could be done using this command
```
bitumctl --testnet --wallet dumpprivkey first_address
```

Next sign raw transaction using private key of first address
```
halfsignedtx=$(bitumctl --testnet --wallet signrawtransaction $rawtxhex "[{\"txid\": \"$utxo_txid\", \"vout\": $utxo_vout, \"tree\": $utxo_tree, \"scriptPubKey\": \"$utxo_spk\", \"redeemScript\": \"$redeem_script\"}]" "[\"first_address_privat_key\"]" | jq -r '.hex')
```

After that it is needed to send hash (halfsignedtx) to second address owner, so he could also sign it
It is needed to get second address private key
```
bitumctl --testnet --wallet dumpprivkey second_address
```

And using it sign hex got by signing of first address private key
```
signedtx=$(bitumctl --testnet --wallet signrawtransaction $halfsignedtx "[{\"txid\": \"$utxo_txid\", \"vout\": $utxo_vout, \"tree\": $utxo_tree, \"scriptPubKey\": \"$utxo_spk\", \"redeemScript\": \"$redeem_script\"}]" "[\"second_address_privat_key\"]" | jq -r '.hex')
```

After that fully signed transaction have to be sended to network for execution
```
bitumctl --testnet --wallet sendrawtransaction $signedtx
```
