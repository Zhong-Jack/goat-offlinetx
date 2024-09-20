package query

import (
	"context"
	"fmt"
	bitcointypes "github.com/goatnetwork/goat/x/bitcoin/types"
	"goat-offlinetx/sdkclient"
)

func QueryBitcoinParams() {
	client := bitcointypes.NewQueryClient(sdkclient.MeClient.GRPCClient)
	response, err := client.Params(context.Background(), &bitcointypes.QueryParamsRequest{})
	if err != nil {
		fmt.Println("query bitcoin params err: ", err)
	}

	fmt.Println("response bitcoin params: ", response.String())
}
