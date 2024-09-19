package query

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	bitcointypes "github.com/goatnetwork/goat/x/bitcoin/types"
	relayertypes "github.com/goatnetwork/goat/x/relayer/types"
	"goat-offlinetx/sdkclient"
)

func QueryGenesis() {
	genesis, err := sdkclient.MeClient.HTTPClient.Genesis(context.Background())
	if err != nil {
		fmt.Println(err)
	}

	var appState map[string]json.RawMessage
	if err := json.Unmarshal(genesis.Genesis.AppState, &appState); err != nil {
		fmt.Errorf("error unmarshalling genesis doc: %s", err)
	}

	var bitcoinState bitcointypes.GenesisState
	if err := sdkclient.MeClient.Cdc.UnmarshalJSON(appState[bitcointypes.ModuleName], &bitcoinState); err != nil {
		fmt.Println("unmarshal error: ", err)
	}
	fmt.Println("bitcoinState: ", bitcoinState)
	fmt.Println("-----------bitcoinState.Pubkey.GetSecp256K1().bytes----------------: ", bitcoinState.Pubkey.GetSecp256K1())
	fmt.Println("-----------bitcoinState.Pubkey.GetSecp256K1().string----------------: ", hex.EncodeToString(bitcoinState.Pubkey.GetSecp256K1()))

	handlePubKey(bitcoinState.Pubkey)

	//var goatState goattypes.GenesisState
	//if err := sdkclient.MeClient.Cdc.UnmarshalJSON(appState[goattypes.ModuleName], &goatState); err != nil {
	//	fmt.Println("unmarshal error: ", err)
	//}
	//fmt.Println("goatState: ", goatState)
	//
	//var relayerState relayertypes.GenesisState
	//if err := sdkclient.MeClient.Cdc.UnmarshalJSON(appState[relayertypes.ModuleName], &relayerState); err != nil {
	//	fmt.Println("unmarshal error: ", err)
	//}
	//fmt.Println("relayerState: ", relayerState)
	//
	//for s, voter := range relayerState.Voters {
	//	fmt.Println("voter address: ", s)
	//	fmt.Println("voter: ", voter)
	//}
}

func handlePubKey(pubKey *relayertypes.PublicKey) {
	bytePubKey := relayertypes.EncodePublicKey(pubKey)

	newPubKey := relayertypes.DecodePublicKey(bytePubKey)
	fmt.Println("newPubKey: ", newPubKey.GetSecp256K1())

	key := hex.EncodeToString(bytePubKey)
	fmt.Println("key: ", key)
	decodeString, err := hex.DecodeString(key)
	if err != nil {
		fmt.Println("decodeString error: ", err)
		return
	}
	newPubKey2 := relayertypes.DecodePublicKey(decodeString)
	fmt.Println("newPubKey2: ", newPubKey2.GetSecp256K1())
}
