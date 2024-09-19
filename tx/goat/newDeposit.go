package goat

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/btcsuite/btcd/wire"
	"github.com/ethereum/go-ethereum/common"
	bitcointypes "github.com/goatnetwork/goat/x/bitcoin/types"
	relayertypes "github.com/goatnetwork/goat/x/relayer/types"
	"github.com/pkg/errors"
	"goat-offlinetx/sdkclient"
)

func NewDepositTx(relayerAddress string) ([]byte, error) {
	if sdkclient.MeClient.HTTPClient == nil {
		return nil, errors.New("client is nil")
	}

	decodeStringRawTX, err := hex.DecodeString("0200000000010149fd0a672a11c5a26cc4b5590d93ae9cd3569898ef0736ffe695a720044c455d0000000000fdffffff0218ddf5050000000022002050a927cae8642c01d76271c28e62a9c887b06824171b431aadfefa784ba8798c00000000000000001a6a184754543070997970c51812dc3a010c7d01b50e0d17dc79c802473044022032f429aa41fb40bbddc814eaf291cb0c1fb2e2c562e9b2deca994844cf728fff02203ad02bf7a5b7ef4f76f6b05ddd3ca8e0f515f82fce012a74ac53ef2140867cd7012102f23e6f933def8455320da0bbe36a432ce8173174eb72c1600743ea22579ebdbc00000000")
	if err != nil {
		return nil, err
	}
	noWitnessTx, err := SerializeNoWitnessTx(decodeStringRawTX)
	if err != nil {
		return nil, err
	}

	address := common.HexToAddress("0x70997970C51812dc3A010C7d01b50e0d17dc79C8").Bytes()
	headers := make(map[uint64][]byte)
	headers[100] = []byte("02894e0a134d6d18c32cd1be11926de8470e01cdc884e3631a4510f0f3ab0f83")
	marshal, err := json.Marshal(headers)
	if err != nil {
		return nil, err
	}

	key := "000383560def84048edefe637d0119a4428dd12a42765a118b2bf77984057633c50e"
	decodeStringKey, err := hex.DecodeString(key)
	if err != nil {
		fmt.Println("decodeString error: ", err)
		return nil, err
	}
	pubKey := relayertypes.DecodePublicKey(decodeStringKey)

	deposits := make([]*bitcointypes.Deposit, 1)
	deposits[0] = &bitcointypes.Deposit{
		Version:           1,
		BlockNumber:       100,
		TxIndex:           1,
		NoWitnessTx:       noWitnessTx,
		OutputIndex:       0,
		IntermediateProof: []byte("02"),
		EvmAddress:        address,
		RelayerPubkey:     pubKey,
	}

	msg := &bitcointypes.MsgNewDeposits{
		Proposer:     "goat1un2gptrl4ecx6v9hwjjz7yw5xf7yyvt6zehxvl",
		BlockHeaders: marshal,
		Deposits:     deposits,
	}

	privKey, err := sdkclient.MeWallet.PrivKey(relayerAddress)
	if err != nil {
		return nil, errors.Wrap(err, "err address private key")
	}
	seq := sdkclient.MeWallet.IncrementNonce(relayerAddress)
	num := sdkclient.MeWallet.AccountNum(relayerAddress)

	tx, err := sdkclient.MeClient.BuildTx(msg, privKey, seq, num)
	if err != nil {
		return nil, errors.Wrap(err, "err BuildTx")
	}

	txBytes, err := sdkclient.MeClient.TxConfig.TxEncoder()(tx)
	if err != nil {
		return nil, errors.Wrap(err, "err TxEncoder")
	}

	return txBytes, nil
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
