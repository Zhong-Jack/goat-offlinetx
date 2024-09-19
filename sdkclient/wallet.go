package sdkclient

import (
	"context"
	"encoding/hex"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/pkg/errors"
)

type Wallet struct {
	addrPrivs map[string]string
	privAddrs map[string]string
	nonces    map[string]*Nonce
	accNums   map[string]uint64
}

func NewWallet(client Client, privs []string) (Wallet, error) {
	var (
		w   Wallet
		err error
	)

	if privs == nil || len(privs) == 0 {
		return Wallet{
			addrPrivs: map[string]string{},
			privAddrs: map[string]string{},
			nonces:    map[string]*Nonce{},
			accNums:   map[string]uint64{},
		}, nil
	}

	if client.HTTPClient == nil {
		return w, errors.New("MeClient is nil")
	}
	//ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	//defer cancel()
	ctx := context.Background()

	w.addrPrivs = make(map[string]string, len(privs))
	w.privAddrs = make(map[string]string, len(privs))
	w.nonces = make(map[string]*Nonce, len(privs))
	w.accNums = make(map[string]uint64, len(privs))
	for _, privKey := range privs {
		var address string
		address, err = GetAddress(privKey)
		if err != nil {
			return Wallet{}, err
		}

		var nonce Nonce
		if nonce, err = NewNonce(ctx, client, address); err != nil {
			err = errors.Wrap(err, "err NewNonce")
			return Wallet{}, err
		}

		var acc types.AccountI
		acc, err = client.Account(ctx, address)
		if err != nil {
			return Wallet{}, err
		}

		w.addrPrivs[address] = privKey
		w.privAddrs[privKey] = address
		w.nonces[address] = &nonce
		w.accNums[address] = acc.GetAccountNumber()
	}
	return w, nil
}

func ImportWallet(privKeyStr string) error {
	if MeClient.HTTPClient == nil {
		return errors.New("MeClient is nil")
	}

	//ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	//defer cancel()
	ctx := context.Background()

	return MeWallet.ImportPrivKey(ctx, MeClient, privKeyStr)
}

func (w *Wallet) ImportPrivKey(ctx context.Context, client Client, privKey string) error {
	var err error
	var address string
	var nonce Nonce
	var acc types.AccountI
	address, err = GetAddress(privKey)
	if err != nil {
		return errors.Wrap(err, "err GetAddress")
	}

	if nonce, err = NewNonce(ctx, client, address); err != nil {
		return errors.Wrap(err, "err NewNonce")
	}

	acc, err = client.Account(ctx, address)
	if err != nil {
		return errors.Wrap(err, "err GetAccount")
	}

	w.addrPrivs[address] = privKey
	w.privAddrs[privKey] = address
	w.nonces[address] = &nonce
	w.accNums[address] = acc.GetAccountNumber()
	return nil
}

func (w *Wallet) PrivKeyStr(address string) string {
	return w.addrPrivs[address]
}

func (w *Wallet) PrivKey(address string) (*secp256k1.PrivKey, error) {
	privKeyStr := w.PrivKeyStr(address)
	if privKeyStr == "" {
		return nil, errors.New("have no private key")
	}
	return GetPrivateKey(privKeyStr)
}

func (w *Wallet) PubKeyStr(privKeyStr string) string {
	return w.privAddrs[privKeyStr]
}

func (w *Wallet) IncrementNonce(address string) uint64 {
	return w.nonces[address].Increment()
}

func (w *Wallet) CurrentNonce(address string) uint64 {
	return w.nonces[address].Current()
}

func (w *Wallet) RecetNonce(address string, nonce uint64) {
	w.nonces[address].Reset(nonce)
}

func (w *Wallet) AccountNum(address string) uint64 {
	return w.accNums[address]
}

func (w *Wallet) Address(privateStr string) (string, error) {
	address := w.privAddrs[privateStr]
	if address == "" {
		return "", errors.Errorf("private key: %s have not address", privateStr)
	}
	return address, nil
}

func (w *Wallet) AccAddr(privateStr string) (sdk.AccAddress, error) {
	address, err := w.Address(privateStr)
	if err != nil {
		return nil, errors.Wrap(err, "err Addr")
	}
	addr, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		return nil, errors.Wrap(err, "err AccAddressFromBech32")
	}
	return addr, nil
}

func GetAddr(privKeyStr string) (sdk.AccAddress, error) {
	privKey, err := GetPrivateKey(privKeyStr)
	if err != nil {
		return nil, err
	}
	return privKey.PubKey().Address().Bytes(), nil
}

func GetAddress(privKeyStr string) (string, error) {
	privKey, err := GetPrivateKey(privKeyStr)
	if err != nil {
		return "", err
	}
	return sdk.AccAddress(privKey.PubKey().Address().Bytes()).String(), nil
}

func GetPrivateKey(privKeyStr string) (*secp256k1.PrivKey, error) {
	priBytes, err := hex.DecodeString(privKeyStr)
	if err != nil {
		return nil, err
	}
	return &secp256k1.PrivKey{Key: priBytes}, nil
}

func CreateRandomAccount() string {
	privKey := secp256k1.GenPrivKey()
	return hex.EncodeToString(privKey.Bytes())
}

func GetPrivateKeyStr(privKey *secp256k1.PrivKey) string {
	return hex.EncodeToString(privKey.Bytes())
}
