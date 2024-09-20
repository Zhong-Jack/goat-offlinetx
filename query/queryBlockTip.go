package query

import (
	"context"
	"fmt"
	bitcointypes "github.com/goatnetwork/goat/x/bitcoin/types"
	"goat-offlinetx/sdkclient"
)

func QueryBlockTip() {
	client := bitcointypes.NewQueryClient(sdkclient.MeClient.GRPCClient)
	response, err := client.BlockTip(context.Background(), &bitcointypes.QueryBlockTipRequest{})
	if err != nil {
		fmt.Println("query bitcoin block tip err: ", err)
	}

	fmt.Println("response block tip: ", response.String())
}
