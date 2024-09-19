package sdkclient

import (
	"context"
	"cosmossdk.io/math"
	"encoding/hex"
	rpchttp "github.com/cometbft/cometbft/rpc/client/http"
	sdkclient "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/std"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	txtypes "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	grouptypes "github.com/cosmos/cosmos-sdk/x/group"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/pkg/errors"
	"goat-offlinetx/config"
	"google.golang.org/grpc"
)

var (
	MeClient Client
	MeWallet Wallet
)

func InitClient(url26657 string, url9090 string) error {
	var err error
	sdkConfig := sdk.GetConfig()
	sdkConfig.SetBech32PrefixForAccount(config.AccountPrefix, config.AccountPrefix+sdk.PrefixPublic)
	sdkConfig.SetBech32PrefixForValidator(config.ValidatorPrefix, config.AccountPrefix+sdk.PrefixPublic)
	sdkConfig.Seal()

	MeClient, err = NewClient(url26657, url9090)
	if err != nil {
		return errors.Wrap(err, "err NewClient")
	}

	MeWallet, _ = NewWallet(MeClient, nil)

	// register custom Denom
	//RegisterDenoms()

	return nil
}

// RegisterDenoms registers token denoms
func RegisterDenoms() {
	err := sdk.RegisterDenom("goat", math.LegacyOneDec())
	if err != nil {
		panic(err)
	}

	err = sdk.RegisterDenom("ugoat", math.LegacyNewDec(6))
	if err != nil {
		panic(err)
	}
}

type Client struct {
	GRPCClient *grpc.ClientConn

	HTTPClient *rpchttp.HTTP
	AuthClient authtypes.QueryClient
	Cdc        *codec.ProtoCodec
	TxConfig   sdkclient.TxConfig
}

func NewClient(url26657 string, url9090 string) (Client, error) {
	var (
		cli = Client{}
		//encCfg = app.MakeEncodingConfig()
		err error
	)

	if cli.HTTPClient, err = rpchttp.New(url26657, "/websocket"); err != nil {
		return cli, err
	}

	if cli.GRPCClient, err = grpc.Dial(url9090, grpc.WithInsecure()); err != nil {
		return cli, err
	}

	amino := codec.NewLegacyAmino()
	std.RegisterLegacyAminoCodec(amino)
	authtypes.RegisterLegacyAminoCodec(amino)
	banktypes.RegisterLegacyAminoCodec(amino)
	grouptypes.RegisterLegacyAminoCodec(amino)
	stakingtypes.RegisterLegacyAminoCodec(amino)
	cli.AuthClient = authtypes.NewQueryClient(cli.GRPCClient)

	//cli.Cdc = codec.NewProtoCodec(encCfg.InterfaceRegistry)

	interfaceRegistry := codectypes.NewInterfaceRegistry()
	cli.Cdc = codec.NewProtoCodec(interfaceRegistry)
	authtypes.RegisterInterfaces(interfaceRegistry)
	std.RegisterInterfaces(interfaceRegistry)
	banktypes.RegisterInterfaces(interfaceRegistry)
	grouptypes.RegisterInterfaces(interfaceRegistry)
	stakingtypes.RegisterInterfaces(interfaceRegistry)

	cli.TxConfig = txtypes.NewTxConfig(cli.Cdc, txtypes.DefaultSignModes)
	return cli, nil
}

func (c Client) LatestBlockHeight(ctx context.Context) (uint64, error) {
	res, err := c.HTTPClient.Block(ctx, nil)
	if err != nil {
		return 0, err
	}
	return uint64(res.Block.Header.Height), nil
}

func (c Client) CountTx(ctx context.Context, height uint64) (int, error) {
	h := int64(height)
	res, err := c.HTTPClient.Block(ctx, &h)
	if err != nil {
		return 0, err
	}
	return len(res.Block.Data.Txs), nil
}

func (c Client) CountPendingTx(ctx context.Context) (int, error) {
	res, err := c.HTTPClient.UnconfirmedTxs(ctx, nil)
	if err != nil {
		return 0, err
	}
	return int(res.Total), nil
}

func (c Client) Nonce(ctx context.Context, address string) (uint64, error) {
	acc, err := c.Account(ctx, address)
	if err != nil {
		return 0, err
	}
	return acc.GetSequence(), nil
}

func (c Client) Account(ctx context.Context, address string) (acc authtypes.AccountI, err error) {
	req := &authtypes.QueryAccountRequest{Address: address}
	res, err := c.AuthClient.Account(ctx, req)
	if err != nil {
		return
	}

	if err = c.Cdc.UnpackAny(res.GetAccount(), &acc); err != nil {
		return
	}

	return
}

func (c Client) Close() {
	c.GRPCClient.Close()
}

func (c *Client) BuildTx(msg sdk.Msg, priv cryptotypes.PrivKey, accSeq uint64, accNum uint64) (authsigning.Tx, error) {
	var txBuilder = c.TxConfig.NewTxBuilder()

	err := txBuilder.SetMsgs(msg)
	if err != nil {
		return nil, err
	}

	// 设置手续费
	fees := sdk.NewCoins(sdk.NewInt64Coin(config.UDenom, 100))
	txBuilder.SetGasLimit(uint64(flags.DefaultGasLimit))
	txBuilder.SetFeeAmount(fees)

	// First round: we gather all the signer infos. We use the "set empty signature" hack to do that.
	if err = txBuilder.SetSignatures(signing.SignatureV2{
		PubKey: priv.PubKey(),
		Data: &signing.SingleSignatureData{
			SignMode:  signing.SignMode(c.TxConfig.SignModeHandler().DefaultMode()),
			Signature: nil,
		},
		Sequence: accSeq,
	}); err != nil {
		return nil, err
	}

	// Second round: all signer infos are set, so each signer can sign.
	signerData := authsigning.SignerData{
		ChainID:       config.ChainID,
		AccountNumber: accNum,
		Sequence:      accSeq,
	}
	sigV2, err := tx.SignWithPrivKey(context.Background(), signing.SignMode(c.TxConfig.SignModeHandler().DefaultMode()), signerData, txBuilder, priv, c.TxConfig, accSeq)
	if err != nil {
		return nil, err
	}
	if err = txBuilder.SetSignatures(sigV2); err != nil {
		return nil, err
	}

	return txBuilder.GetTx(), nil
}

func AccAddressFromPrivString(privStr string) (string, error) {
	priBytes, err := hex.DecodeString(privStr)
	if err != nil {
		return "", err
	}
	priv := &secp256k1.PrivKey{Key: priBytes}
	accAddress := sdk.AccAddress(priv.PubKey().Address().Bytes()).String()

	return accAddress, nil
}
