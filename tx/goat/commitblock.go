package goat

import (
	bitcointypes "github.com/goatnetwork/goat/x/bitcoin/types"
	relayertypes "github.com/goatnetwork/goat/x/relayer/types"
	"github.com/pkg/errors"
	"goat-offlinetx/sdkclient"
)

func CommitBlock(relayerAddress string) ([]byte, error) {
	if sdkclient.MeClient.HTTPClient == nil {
		return nil, errors.New("client is nil")
	}

	msg := &bitcointypes.MsgNewBlockHashes{
		Proposer:         "goat1un2gptrl4ecx6v9hwjjz7yw5xf7yyvt6zehxvl",
		StartBlockNumber: 100,
		BlockHash: [][]byte{
			[]byte("02894e0a134d6d18c32cd1be11926de8470e01cdc884e3631a4510f0f3ab0f83"),
		},
		Vote: &relayertypes.Votes{
			Signature: []byte("signature"),
		},
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
