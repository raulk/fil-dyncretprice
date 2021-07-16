package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/big"
	"github.com/filecoin-project/lotus/build"
	"github.com/ipfs/go-cid"
	"github.com/libp2p/go-libp2p-core/peer"
)

type Ask struct {
	PricePerByte            abi.TokenAmount
	UnsealPrice             abi.TokenAmount
	PaymentInterval         uint64
	PaymentIntervalIncrease uint64
}

type PricingInput struct {
	// PayloadCID is the cid of the payload to retrieve.
	PayloadCID cid.Cid
	// PieceCID is the cid of the Piece from which the Payload will be retrieved.
	PieceCID cid.Cid
	// PieceSize is the size of the Piece from which the payload will be retrieved.
	PieceSize abi.UnpaddedPieceSize
	// Client is the peerID of the retrieval client.
	Client peer.ID
	// VerifiedDeal is true if there exists a verified storage deal for the PayloadCID.
	VerifiedDeal bool
	// Unsealed is true if there exists an unsealed sector from which we can retrieve the given payload.
	Unsealed bool
	// CurrentAsk is the ask configured in the ask-store via the `lotus retrieval-deals set-ask` CLI command.
	CurrentAsk Ask
}

var AsksByPeer = map[string]Ask{
	"peerid1": {
		PricePerByte: abi.NewTokenAmount(int64(0.0003 * float64(build.FilecoinPrecision))),
		UnsealPrice:  big.Zero(),
	},
	// more...
}

func main() {
	var (
		// decode/encode stdin/stdout as json lines.
		dec = json.NewDecoder(os.Stdin)
		enc = json.NewEncoder(os.Stdout)
	)

	// parse the pricing input.
	var in PricingInput
	if err := dec.Decode(&in); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to decode input as PricingInput json: %s", err)
		os.Exit(1)
		return
	}

	// do we have a special ask for this peer?
	ret, ok := AsksByPeer[in.Client.String()]
	if !ok {
		// we don't; fall back to the default configured ask, which is provided
		// to us in the input message.
		ret = in.CurrentAsk
	}

	// write the response.
	if err := enc.Encode(ret); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to encode output as Ask json: %s\n", err)
		_, _ = fmt.Fprintf(os.Stderr, "computed Ask was: %+v\n", ret)
		os.Exit(1)
		return
	}
}
