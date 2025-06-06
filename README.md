## Blockwatch TzGo - Tezos Go SDK

TzGo is [Blockwatch](https://blockwatch.cc)'s low-level Tezos Go SDK for reliable, high-performance applications. This SDK is free to use in commercial and non-commercial projects with a permissive license. Blockwatch is committed to keeping interfaces stable, providing long-term support, and updating TzGo on a regular basis to stay compliant with the most recent Tezos network protocol.

TzGo's main focus is on **correctness**, **stability**, and **compliance** with Tezos mainnet. It supports binary and JSON encodings for all Tezos types including Micheline smart contract data and all transaction formats. It's an ideal fit for high-performance applications that read from and write to the Tezos blockchain.

Current Tezos protocol support in TzGo

- Quebec v021
- ParisC v020
- Paris v019
- Oxford v018
- Nairobi v017
- Mumbai v016
- Lima v015
- Kathmandu v014
- Jakarta v013
- Ithaca v012
- Hangzhou v011
- Granada v010
- Florence v009
- Edo v008
- Delphi v007
- Carthage v006
- Babylon v005
- Athens v004
- Alpha v001-v003

### SDK features

TzGo contains a full set of features to read, monitor, decode, translate, analyze and debug data from the Tezos blockchain, in particular from Tezos smart contracts:

- a low-level **Types library** `tzgo/tezos` to handle hashes, addresses, keys, signatures other types found on-chain
- a powerful **Micheline library** `tzgo/micheline` to decode and translate data found in smart contract calls, storage, and bigmaps
- an **RPC library** `tzgo/rpc` for accessing the Tezos Node RPC
- a **Codec library** `tzgo/codec` to construct and serialize operations
- a **Contract library** `tzgo/contract` for smart contract calls and tokens
- a **Signer library** `tzgo/signer` to sign transactions local or remote
- helpers like an efficient base58 en/decoder, hash maps, etc
- a **Code generator** [TzGen](https://github.com/blockwatch-cc/tzgo/tree/master/cmd/tzgen) to produce pure Go clients for smart contract interfaces
- an **Automation Tool** [TzCompose](https://github.com/blockwatch-cc/tzgo/tree/master/cmd/tzcompose) to setup test cases and deploy complex contract ecosystems

### TzGo Compatibility

TzGo's RPC package attempts to be compatible with all protocols so that reading historic block data is always supported. Binary transaction encoding and signing support is limited to the most recent protocol.

We attempt to upgrade TzGo whenever new protocols are proposed and will add new protocol features as soon as practically feasible and as demand for such features exists. For example, we don't fully Sapling and BLS signatures yet, but may add support in the future.

### Usage

```sh
go get -u github.com/trilitech/tzgo
```

Then import, using

```go
import (
	"github.com/trilitech/tzgo/codec"
	"github.com/trilitech/tzgo/tezos"
	"github.com/trilitech/tzgo/micheline"
	"github.com/trilitech/tzgo/rpc"
	"github.com/trilitech/tzgo/wallet"
)
```

### Micheline Support

Tezos uses [Micheline](https://tezos.gitlab.io/shell/micheline.html) for encoding smart contract data and code. The positive is that Micheline is strongly typed, the downside is that it's complex and has a few ambiguities that make it hard to use. TzGo contains a library that lets you decode, analyze and construct compliant Micheline data structures from Go.

Micheline uses basic **primitives** for encoding types and values. These primitives can be expressed in JSON and binary format and TzGo can translate between them efficiently. Micheline also supports type **annotations** which are used by high-level languages to express complex data types like records and their field names.

TzGo defines a basic `Prim` data type to work with Micheline primitives directly:

```go
type Prim struct {
	Type      PrimType // primitive type
	OpCode    OpCode   // primitive opcode (invalid on sequences, strings, bytes, int)
	Args      []Prim   // optional nested arguments
	Anno      []string // optional type annotations
	Int       *big.Int // decoded value when Prim is an int
	String    string   // decoded value when Prim is a string
	Bytes     []byte   // decoded value when Prim is a byte sequence
	WasPacked bool     // true when content has been unpacked
}
```

Since Micheline value encoding is quite verbose and can be ambiguous, TzGo supports **unfolding** of raw Micheline using the following TzGo wrapper types and their helper functions like `Map()`, `GetInt64()`, `GetAddress()`:

- `Type` is a TzGo wrapper for simple or complex primitives which contain annotated type info
- `Value` is a TzGo wrapper for simple or complex primitives representing Micheline values in combination with their Type
- `Key` is a TzGo wrapper for special comparable values that are used as maps or bigmap keys

Sometimes Micheline values have been packed into byte sequences with the Michelson PACK instruction and it is desirable to unpack them before processing (e.g. to retrieve UFT8 strings or nested records). TzGo supports `Unpack()` and `UnpackAll()` functions on primitives and values and also detects the internal type of packed data which is necessary for unfolding.


### Examples

Below are a few examples showing how to use TzGo to easily access Tezos data in your application.

#### Parsing an address

To parse/decode an address and output its components you can do the following:

```go
import "github.com/trilitech/tzgo/tezos"

// parse and panic if invalid
addr := tezos.MustParseAddress("tz3RDC3Jdn4j15J7bBHZd29EUee9gVB1CxD9")

// parse and return error if invalid
addr, err := tezos.ParseAddress("tz3RDC3Jdn4j15J7bBHZd29EUee9gVB1CxD9")
if err != nil {
	fmt.Printf("Invalid address: %v\n", err)
}

// Do smth with the address
fmt.Printf("Address type = %s\n", addr.Type)
fmt.Printf("Address bytes = %x\n", addr.Hash)

```

See [examples/addr.go](https://github.com/blockwatch-cc/tzgo/blob/master/examples/addr/main.go) for more.

#### Monitoring for new blocks

A Tezos node can notify applications when new blocks are attached to the chain. The Tezos RPC calls this monitor and technically it's a long-poll implementation. Here's how to use this feature in TzGo:

```go
import "github.com/trilitech/tzgo/rpc"

// init SDK client
c, _ := rpc.NewClient("https://rpc.tzstats.com", nil)

// create block header monitor
mon := rpc.NewBlockHeaderMonitor()
defer mon.Close()

// all SDK functions take a context, here we just use a dummy
ctx := context.TODO()

// register the block monitor with our client
if err := c.MonitorBlockHeader(ctx, mon); err != nil {
	log.Fatalln(err)
}

// wait for new block headers
for {
	head, err := mon.Recv(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	// do smth with the block header
	fmt.Printf("New block %s\n", head.Hash)
}

```

#### Fetch and decode contract storage

```go
import (
	"github.com/trilitech/tzgo/micheline"
	"github.com/trilitech/tzgo/rpc"
	"github.com/trilitech/tzgo/tezos"
)

// we use the Baker Registry on mainnet as example
addr := tezos.MustParseAddress("KT1ChNsEFxwyCbJyWGSL3KdjeXE28AY1Kaog")

// init RPC client
c, _ := rpc.NewClient("https://rpc.tzstats.com", nil)

// fetch the contract's script and most recent storage
script, _ := c.GetContractScript(ctx, addr)

// unfold Micheline storage into human-readable form
val := micheline.NewValue(script.StorageType(), script.Storage)
m, _ := val.Map()
buf, _ := json.MarshalIndent(m, "", "  ")
fmt.Println(string(buf))
```

#### Decode contract call parameters

```go
import (
	"github.com/trilitech/tzgo/micheline"
	"github.com/trilitech/tzgo/rpc"
	"github.com/trilitech/tzgo/tezos"
)

// init RPC client
c, _ := rpc.NewClient("https://rpc.tzstats.com", nil)

// assuming you have this transaction
tx := block.Operations[3][0].Contents[0].(*rpc.TransactionOp)

// load the contract's script for type info
script, err := c.GetContractScript(ctx, tx.Destination)

// unwrap params for nested entrypoints
ep, param, err := tx.Parameters.MapEntrypoint(script.ParamType())

// convert Micheline param data into human-readable form
val := micheline.NewValue(ep.Type(), param)

// e.g. access individual nested fields using value helpers
from, ok := val.GetAddress("transfer.from")
```

#### Use TzGo's Value API

Micheline type and value trees are verbose and can be ambiguous due to comb-pair optimizations. If you don't know or don't care about what that even means, you may want to use the `Value` API which helps you translate Micheline into human readable form.

There are multiple options to access decoded data:

```go
// 1/
// decode into a Go type tree using `Map()` which produces
// - `map[string]interface{}` for records and maps
// - `[]interface{}` for lists and sets
// - `time.Time` for time values
// - `string` stringified numbers for Nat, Int, Mutez
// - `bool` for Booleans
// - `tezos.Address` for any string or bytes sequence that contains an address
// - `tezos.Key` for keys
// - `tezos.Signature` for signatures
// - `tezos.ChainIdHash` for chain ids
// - hex string for untyped bytes
// - opcode names for any Michelson opcode
// - opcode sequences for lambdas
m, err := val.Map()
buf, err := json.MarshalIndent(m, "", "  ")
fmt.Println(string(buf))

// when you know the concrete type you can cast directly
fmt.Println("Value=", m.(map[string]interface{})["transfer"].(map[string]interface{})["value"])

// 2/
// access individual nested fields using value helpers (`ok` is true when the field
// exists and has the correct type; helpers exist for
//   GetString() (string, bool)
//   GetBytes() ([]byte, bool)
//   GetInt64() (int64, bool)
//   GetBig() (*big.Int, bool)
//   GetBool() (bool, bool)
//   GetTime() (time.Time, bool)
//   GetAddress() (tezos.Address, bool)
//   GetKey() (tezos.Key, bool)
//   GetSignature() (tezos.Signature, bool)
from, ok := val.GetAddress("transfer.from")

// 3/
// unmarshal the decoded Micheline parameters into a json-tagged Go struct
type FA12Transfer struct {
    From  tezos.Address `json:"from"`
    To    tezos.Address `json:"to"`
    Value int64         `json:"value,string"`
}

// FA1.2 params are nested, so we need an extra wrapper
type FA12TransferParams struct {
	Transfer FA12Transfer `json:"transfer"`
}

var transfer FA12TransferParams
err := val.Unmarshal(&transfer)
```

#### List a contract's bigmaps

```go
import (
	"context"

	"github.com/trilitech/tzgo/rpc"
	"github.com/trilitech/tzgo/tezos"
)

// we use the hic et nunc NFT market on mainnet as example
addr := tezos.MustParseAddress("KT1Hkg5qeNhfwpKW4fXvq7HGZB9z2EnmCCA9")

// init RPC client
c, _ := rpc.NewClient("https://rpc.tzpro.io", nil)
ctx := context.TODO()

// fetch the contract's script and most recent storage
script, _ := c.GetContractScript(ctx, addr)

// bigmap pointers as named map[string]int64 (names from type annotations)
// Note that if a bigmap is anonymous, e.g. in a list, a temporary name will
// be returned here
named := script.Bigmaps()

// bigmap pointers as []int64
ids := []int64{}
for _, v := range named {
	ids = append(ids, v)
}
```

#### Fetch and decode bigmap values

```go
// init RPC client
c, _ := rpc.NewClient("https://rpc.tzstats.com", nil)

// load bigmap type info (use the Baker Registry on mainnet as example)
biginfo, _ := c.GetBigmapInfo(ctx, 17)

// list all bigmap keys
bigkeys, _ := c.GetBigmapKeys(ctx, 17)

// visit each value
for _, key := range bigkeys {
	bigval, _ := c.GetBigmapValue(ctx, 17, key)

	// unfold Micheline type into human readable form
	val := micheline.NewValue(micheline.NewType(biginfo.ValueType), bigval)
	m, _ := val.Map()
	buf, _ := json.MarshalIndent(m, "", "  ")
	fmt.Println(string(buf))
}
```

#### Custom RPC client configuration

TzGo's `rpc.NewClient()` function takes an optional Go `http.Client` as parameter which you can configure before or after passing it to the library. The example below shows how to set custom timeouts and disable TLS certificate checks (not recommended in production, but useful if you use self-signed certificates during testing).


```go
import (
	"crypto/tls"
	"log"
	"net"
	"net/http"

	"github.com/trilitech/tzgo/rpc"
)


func main() {
	hc := &http.Client{
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout:   2 * time.Second,
				KeepAlive: 180 * time.Second,
			}).Dial,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			}
		}
	}

	c, err := rpc.NewClient("https://my-private-node.local:8732", hc)
	if err != nil {
		log.Fatalln(err)
	}
}
```


## License

The MIT License (MIT) Copyright (c) 2020-2024 Blockwatch Data Inc.

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is furnished
to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
