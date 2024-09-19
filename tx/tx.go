package tx

import (
	"context"
	"fmt"
	ctypes "github.com/cometbft/cometbft/rpc/core/types"
	sdktx "github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/pkg/errors"
	"goat-offlinetx/sdkclient"
)

/*
SendTx 发送单笔交易
@Description:
@param txBytes
@return error
*/
func SendTx(txBytes []byte) error {
	if sdkclient.MeClient.HTTPClient == nil {
		return errors.New("client is nil")
	}

	//ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	//defer cancel()
	ctx := context.Background()

	txClient := sdktx.NewServiceClient(sdkclient.MeClient.GRPCClient)

	grpcRes, err := txClient.BroadcastTx(ctx, &sdktx.BroadcastTxRequest{
		Mode:    sdktx.BroadcastMode_BROADCAST_MODE_SYNC,
		TxBytes: txBytes},
	)
	if err != nil {
		fmt.Println("BroadcastTx is err:", err)
		return err
	}

	fmt.Println("打印 tx response info:", grpcRes)

	return nil
}

/*
SendBatchTX 交易批处理
@Description:
@param txs 多笔交易信息
@return error
*/
func SendBatchTX(txs [][]byte) error {
	if sdkclient.MeClient.HTTPClient == nil {
		return errors.New("client is nil")
	}

	//ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	//defer cancel()
	ctx := context.Background()

	batch := sdkclient.MeClient.HTTPClient.NewBatch()
	// Broadcast the transaction and wait for it to commit (rather use c.BroadcastTxSync though in production).
	for _, tx := range txs {
		_, err := batch.BroadcastTxSync(ctx, tx)
		if err != nil {
			return errors.Wrap(err, "BroadcastTxSync err")
		}
	}

	// Send the batch of more than 2 transactions
	results, err := batch.Send(ctx)
	if err != nil {
		return errors.Wrap(err, "err Send")
	}
	for _, res := range results {
		txResult, _ := res.(*ctypes.ResultBroadcastTx)
		fmt.Println("print tx response info:", txResult)
	}

	return nil
}
