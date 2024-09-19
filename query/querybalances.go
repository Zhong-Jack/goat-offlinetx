package query

import (
	"context"
	"fmt"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"goat-offlinetx/config"
	"goat-offlinetx/sdkclient"
)

func QueryBalances(addr string) {
	bankRequest := banktypes.QueryBalanceRequest{
		Address: addr,
		Denom:   config.UDenom,
	}

	client := banktypes.NewQueryClient(sdkclient.MeClient.GRPCClient)
	response, err := client.Balance(context.Background(), &bankRequest)
	if err != nil {
		fmt.Println("query balances err：", err)

	}

	fmt.Println("query balances response: ", response)
}
