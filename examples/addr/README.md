## Read Tezos addresses and public keys

This example shows how Tezos addresses and keys can be detected and decoded.

### Code Example

```go
// pubkey
key, err = mavryk.ParseKey("edpk....")

// addresss
addr, err = mavryk.ParseAddress("mv1...")
```

### Usage

```sh
// decode address types
go run . mv1CitY5z26GMMmEHFDVyG7ksMv3XPpqEsbS
go run . mv2Nbn1Vg9s93REgZ3QWjgHmgR8yBQoHuvsh
go run . mv3M6rekrGAYypmia876dVkoC3N9Uuj7xVXh
go run . KT1GyeRktoGPEKsWpchWguyy8FAf3aNHkw2T

// decode pubkeys
go run . edpkucde3WUTR2s6KgDBwvR7NiezGyHNj1aGz6WrJg6SeZWHNjDA8N
go run . sppk7aAV5AjmQPcph9SrrKBBeFwj15kMvnByjbvb9mqsTMgUm1ZoHxK
go run . p2pk68CeMSnZ8MhrW6zCJzGfS2VTsFUKK5GwB7Hem3UUuyQH2kHHeij

// create blinded address from address and secret
go run . tz1T1rRqmAk4XtGadNJuNpq8dUdWqLv2Gtq4 06da1e038224114366831e47aee7f128f4675311
```