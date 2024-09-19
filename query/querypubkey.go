package query

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/base58"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	bitcointypes "github.com/goatnetwork/goat/x/bitcoin/types"
	"goat-offlinetx/sdkclient"
	"golang.org/x/crypto/ripemd160"
)

func QueryPubKey() {
	client := bitcointypes.NewQueryClient(sdkclient.MeClient.GRPCClient)
	response, err := client.Pubkey(context.Background(), &bitcointypes.QueryPubkeyRequest{})
	if err != nil {
		fmt.Println("query goat pubkey err：", err)
	}

	pubKey := response.PublicKey.GetSecp256K1()

	// 将公钥字节数组转换为十六进制字符串
	pubKeyHex := hex.EncodeToString(pubKey)
	fmt.Println("公钥的十六进制字符串:", pubKeyHex)

	// 将十六进制字符串转换为 Base64 字符串
	pubKeyBase64 := base64.StdEncoding.EncodeToString([]byte(pubKeyHex))
	fmt.Println("公钥的 Base64 字符串:", pubKeyBase64)

	fmt.Println("pubkey: ", pubKey)

	// 计算 SHA-256 哈希
	sha256Hash := sha256.Sum256(pubKey)

	// 计算 RIPEMD-160 哈希
	ripemd160Hash := ripemd160.New()
	ripemd160Hash.Write(sha256Hash[:])
	publicKeyHash := ripemd160Hash.Sum(nil)

	// 添加网络前缀（主网地址前缀为 0x00）
	versionedPayload := append([]byte{0x00}, publicKeyHash...)

	// 计算校验和
	checksum := sha256.Sum256(versionedPayload)
	checksum = sha256.Sum256(checksum[:])
	checksum2 := checksum[:4] // 取前4个字节

	// 构建完整的地址
	addressBytes := append(versionedPayload, checksum2...)

	// 编码为 Base58Check
	address := base58.Encode(addressBytes)

	// 输出比特币地址
	fmt.Println("比特币地址:", address)
}

func QueryPubKey2() {
	client := bitcointypes.NewQueryClient(sdkclient.MeClient.GRPCClient)
	response, err := client.Pubkey(context.Background(), &bitcointypes.QueryPubkeyRequest{})
	if err != nil {
		fmt.Println("query goat pubkey err：", err)
	}

	pubKey := response.PublicKey.GetSecp256K1()

	fmt.Println("pubkey: ", pubKey)

	var posPubkey *btcec.PublicKey

	posPubkey, _ = btcec.ParsePubKey(pubKey)

	evmAddress, _ := hex.DecodeString("0xb43662708D58854EE74cE844C79732A32fF8C195")

	redeemScript, err := txscript.NewScriptBuilder().
		AddData(evmAddress[:]).
		AddOp(txscript.OP_DROP).
		AddData(posPubkey.SerializeCompressed()).
		AddOp(txscript.OP_CHECKSIG).Script()
	if err != nil {
		fmt.Println("build redeem script err:", err)
	}

	witnessProg := sha256.Sum256(redeemScript)
	network := &chaincfg.MainNetParams
	depositAddress, err := btcutil.NewAddressWitnessScriptHash(witnessProg[:], network)
	if err != nil {
		fmt.Println("new address err:", err)
	}
	// bc1ptc5fh8zc06j427c92yk8uy3aytey0r9yvqe8aagrjccfrw6d3cjshvhj4n
	fmt.Println("depositAddress: ", depositAddress)
}

func QueryPubKey3() {
	client := bitcointypes.NewQueryClient(sdkclient.MeClient.GRPCClient)
	response, err := client.Pubkey(context.Background(), &bitcointypes.QueryPubkeyRequest{})
	if err != nil {
		fmt.Println("query goat pubkey err：", err)
	}

	pubKey := response.PublicKey.GetSecp256K1()

	// 将公钥字节数组转换为十六进制字符串
	pubKeyHex := hex.EncodeToString(pubKey)
	fmt.Println("公钥的十六进制字符串:", pubKeyHex)

	network := &chaincfg.MainNetParams
	p2wpkh, err := btcutil.NewAddressWitnessPubKeyHash(btcutil.Hash160(pubKey), network)
	if err != nil {
		fmt.Println("new address err:", err)
	}

	fmt.Println("btc address: ", p2wpkh.EncodeAddress())
}

func QueryDepositAddress() {
	client := bitcointypes.NewQueryClient(sdkclient.MeClient.GRPCClient)
	response, err := client.DepositAddress(context.Background(), &bitcointypes.QueryDepositAddress{
		Version:    1,
		EvmAddress: "0x454279b3519cE722a3Bf04983cE41793693BB758",
	})
	if err != nil {
		fmt.Println("query goat deposit address err：", err)
	}

	fmt.Println("deposit address: ", response.Address)
}
