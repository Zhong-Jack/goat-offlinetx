package query

import (
	"context"
	"fmt"
	relayertypes "github.com/goatnetwork/goat/x/relayer/types"
	"goat-offlinetx/sdkclient"
)

func QueryRelayer() {
	client := relayertypes.NewQueryClient(sdkclient.MeClient.GRPCClient)
	response, err := client.Relayer(context.Background(), &relayertypes.QueryRelayerRequest{})
	if err != nil {
		fmt.Println("query relayer err: ", err)
	}

	fmt.Println("response sequence: ", response.Sequence)

	fmt.Println("response relayer: ", response.Relayer)

	voters := response.Relayer.Voters
	for _, voter := range voters {
		fmt.Println("voter: ", voter)
	}
}
