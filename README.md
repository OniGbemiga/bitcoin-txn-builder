# simple-bitcoin-transaction-builder

## Prerequisites
* Have an understanding of bitcoin
* Have GO set up on your development environment
* Have Bitcoin core set up and running on your development environment
* Have access to a terminal for CLI operations

## Getting Started

Clone this GitHub repo to your development environment and run the command below
in the root directory of the folder to download the necessary packages:

```
go mod download
```

## CLI Actions
These actions should all the performed in a terminal opened in the root directory of the cloned repo

1. Generate Redeem Script Hex

```
go run cmd/main.go -task=redeem -preimage="Btrust Builders"
```
2. Generate Address from the Redeem Script Hex

```
go run cmd/main.go -task=derive -preimage="Btrust Builders"
```

3. Construct a transaction that sends bitcoin to the generated address

```
go run cmd/main.go -task=transaction -preimage="Btrust Builders" -dest="destination_address" private-key="your_private_key" -amount=amount
```

4. Construct a transaction that spends from the transaction above

```
go run cmd/main.go -task=spending -preimage="Btrust Builders" private-key="your_private_key", txHexPtr="txHex from number 3 above"
```
