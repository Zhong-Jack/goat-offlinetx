package main

import (
	"fmt"
	"goat-offlinetx/config"
	"goat-offlinetx/query"
	"goat-offlinetx/sdkclient"
)

func main() {
	err := sdkclient.InitClient(config.DefaultRPCURI, config.DefaultGRPCURI)
	if err != nil {
		fmt.Println("InitClient err: ", err)
		return
	}

	err = sdkclient.ImportWallet(config.RelayerPriKey)
	if err != nil {
		fmt.Println("ImportWallet err: ", err)
		return
	}

	////--------------------------------------提交区块交易------------------------------------------------//
	//// 签名地址
	//relayerAddr, _ := sdkclient.MeWallet.Address(config.RelayerPriKey)
	//fmt.Println("签名者地址:", relayerAddr)
	//
	//// 构建交易
	//txBytes, err := goat.CommitBlock(relayerAddr)
	//if err != nil {
	//	fmt.Println("构建CommitBlock交易失败：", err)
	//	return
	//}
	//
	//// 发送交易
	//err = tx.SendTx(txBytes)
	//if err != nil {
	//	fmt.Println("发送交易失败: ", err)
	//	return
	//}

	////--------------------------------------提交deposit交易------------------------------------------------//
	//// 签名地址
	//relayerAddr, _ := sdkclient.MeWallet.Address(config.RelayerPriKey)
	//fmt.Println("--------------relayerAddr----------: ", relayerAddr)
	//
	//// 构建交易
	//txBytes, err := goat.NewDepositTx(relayerAddr)
	//if err != nil {
	//	fmt.Println("构建NewDepositTx失败：", err)
	//	return
	//}
	//
	//// 发送交易
	//err = tx.SendTx(txBytes)
	//if err != nil {
	//	fmt.Println("发送交易失败: ", err)
	//	return
	//}

	//--------------------------------------查询pubkey---------------------------------------------//
	//query.QueryPubKey()
	//query.QueryPubKey2()
	//query.QueryPubKey3()

	//--------------------------------------查询deposit address------------------------------------//
	//query.QueryDepositAddress()

	//--------------------------------------查询genesis--------------------------------------------//
	//query.QueryGenesis()

	//--------------------------------------查询relayer--------------------------------------------//
	//query.QueryRelayer()

	//--------------------------------------查询blockTip--------------------------------------------//
	//query.QueryBlockTip()

	//--------------------------------------查询deposit--------------------------------------------//
	//query.QueryDeposit()

	//--------------------------------------查询bitcoin params--------------------------------------//
	query.QueryBitcoinParams()
}
