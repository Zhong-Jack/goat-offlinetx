package query

import (
	"bytes"
	"context"
	"encoding/hex"

	"fmt"
	"github.com/btcsuite/btcd/wire"
	goatcrypto "github.com/goatnetwork/goat/pkg/crypto"
	bitcointypes "github.com/goatnetwork/goat/x/bitcoin/types"
	"goat-offlinetx/sdkclient"
)

func QueryDeposit() {
	decodeString, err := hex.DecodeString("0200000000010125296d6dd6f9b5c1a10cb1bd2ec0dd55feaabb3f5ff4b8ffd7133e89a591f5b50000000000fdffffff0218ddf505000000001600149759ed6aae6ade43ae6628a943a39974cd21c5df00000000000000001a6a184754543070997970c51812dc3a010c7d01b50e0d17dc79c802473044022011aaf72ac782c780c68190c136be89072cf1f45964bff2047440c4b8e37b225b02206a06abd50de5f7caf39dc01aa7a99b4cfee65e7466980e2ca4a95d7a5a3222e6012102f23e6f933def8455320da0bbe36a432ce8173174eb72c1600743ea22579ebdbc00000000")
	if err != nil {
		fmt.Sprintf("DecodeString err: %v", err)
	}

	noWitnessTx, err := SerializeNoWitnessTx(decodeString)
	if err != nil {
		fmt.Sprintf("SerializeNoWitnessTx err: %v", err)
	}
	txid := goatcrypto.DoubleSHA256Sum(noWitnessTx)

	// txhash: 601b0f071d85cb53aa87b82db7667eca3126b450bf0fa9994948627ee054a8c1
	newTxid := hex.EncodeToString(txid)

	client := bitcointypes.NewQueryClient(sdkclient.MeClient.GRPCClient)
	response, err := client.HasDeposited(context.Background(), &bitcointypes.QueryHasDeposited{
		Txid:  newTxid,
		Txout: 0,
	})
	if err != nil {
		fmt.Println("query deposit err: ", err)
	}

	fmt.Println("response deposit: ", response.String())
}

func SerializeNoWitnessTx(rawTransaction []byte) ([]byte, error) {
	// Parse the raw transaction
	rawTx := wire.NewMsgTx(wire.TxVersion)
	err := rawTx.Deserialize(bytes.NewReader(rawTransaction))
	if err != nil {
		return nil, fmt.Errorf("failed to parse raw transaction: %v", err)
	}

	// Create a new transaction without witness data
	noWitnessTx := wire.NewMsgTx(rawTx.Version)

	// Copy transaction inputs, excluding witness data
	for _, txIn := range rawTx.TxIn {
		newTxIn := wire.NewTxIn(&txIn.PreviousOutPoint, nil, nil)
		newTxIn.Sequence = txIn.Sequence
		noWitnessTx.AddTxIn(newTxIn)
	}

	// Copy transaction outputs
	for _, txOut := range rawTx.TxOut {
		noWitnessTx.AddTxOut(txOut)
	}

	// Set lock time
	noWitnessTx.LockTime = rawTx.LockTime

	// Serialize the transaction without witness data
	var buf bytes.Buffer
	err = noWitnessTx.Serialize(&buf)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize transaction without witness data: %v", err)
	}

	return buf.Bytes(), nil
}
