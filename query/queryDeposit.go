package query

import (
	"context"
	"fmt"
	bitcointypes "github.com/goatnetwork/goat/x/bitcoin/types"
	"goat-offlinetx/sdkclient"
)

func QueryDeposit() {
	client := bitcointypes.NewQueryClient(sdkclient.MeClient.GRPCClient)
	response, err := client.HasDeposited(context.Background(), &bitcointypes.QueryHasDeposited{
		Txid:  "601b0f071d85cb53aa87b82db7667eca3126b450bf0fa9994948627ee054a8c1",
		Txout: 0,
	})
	if err != nil {
		fmt.Println("query deposit err: ", err)
	}

	fmt.Println("response deposit: ", response)
}
