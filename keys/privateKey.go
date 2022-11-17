package keys

import (
	"bytes"
	"crypto/ecdsa"
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethtweet/ethtweet/global"
	"github.com/ethtweet/ethtweet/logs"

	cryptoEth "github.com/ethereum/go-ethereum/crypto"
	keystore "github.com/ipfs/go-ipfs-keystore"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/mr-tron/base58"
)

type PrivateKey struct {
	EthPrivate    *ecdsa.PrivateKey
	LibP2pPrivate crypto.PrivKey
}

func NewPrivateKey() (*PrivateKey, error) {
	priKey, _, err := crypto.GenerateKeyPair(
		crypto.Secp256k1,
		-1,
	)
	if err != nil {
		return nil, err
	}
	ethPri, err := global.LibP2pPriToEthPri(priKey)
	if err != nil {
		return nil, err
	}
	return &PrivateKey{
		EthPrivate:    ethPri,
		LibP2pPrivate: priKey,
	}, nil
}

func NewPrivateKeyByLibP2pPri(priKey crypto.PrivKey) (*PrivateKey, error) {
	ethPri, err := global.LibP2pPriToEthPri(priKey)
	if err != nil {
		return nil, err
	}
	return &PrivateKey{
		EthPrivate:    ethPri,
		LibP2pPrivate: priKey,
	}, nil
}

func NewPrivateKeyByBase58(keyBase string) (*PrivateKey, error) {
	kb, err := base58.Decode(keyBase)
	if err != nil {
		return nil, err
	}
	pri, err := crypto.UnmarshalPrivateKey(kb)
	if err != nil {
		return nil, err
	}
	return NewPrivateKeyByLibP2pPri(pri)
}

func NewPrivateKeyByEthPri(priKey *ecdsa.PrivateKey) (*PrivateKey, error) {
	libP2pPri, err := global.EthPriToLibP2pPri(priKey)
	if err != nil {
		return nil, err
	}
	return &PrivateKey{
		EthPrivate:    priKey,
		LibP2pPrivate: libP2pPri,
	}, nil
}

func (pri *PrivateKey) EncodePublic() []byte {
	return cryptoEth.FromECDSAPub(&pri.EthPrivate.PublicKey)
}

func (pri *PrivateKey) Encode58Public() string {
	return base58.Encode(pri.EncodePublic())
}

func (pri *PrivateKey) PutStore(dir, keyName string) error {
	ks, err := keystore.NewFSKeystore(dir)
	if err != nil {
		return err
	}
	return ks.Put(keyName, pri.LibP2pPrivate)
}

func (pri *PrivateKey) GetEthAddress() common.Address {
	return cryptoEth.PubkeyToAddress(pri.EthPrivate.PublicKey)
}

func (pri *PrivateKey) Sign(msg string) (string, error) {
	s, err := cryptoEth.Sign(global.EthSignHash(msg), pri.EthPrivate)
	if err != nil {
		return "", err
	}
	s[64] += 27
	return hexutil.Encode(s), nil
}

func VerifySignature(publicKey, sign, msg string) bool {
	pubBytes, err := base58.Decode(publicKey)
	if err != nil {
		return false
	}
	recoveredPub, err := FetchPubKeyBySignMsg(sign, msg)
	if err != nil {
		logs.PrintErr(err)
		return false
	}
	if bytes.Equal(cryptoEth.FromECDSAPub(recoveredPub), pubBytes) {
		return true
	}
	return false
}

func VerifySignatureByAddress(address, sign, msg string) bool {
	recoveredPub, err := FetchPubKeyBySignMsg(sign, msg)
	if err != nil {
		logs.PrintErr(err)
		return false
	}
	pAddress := PubKeyToAddress(recoveredPub)
	if pAddress.String() == address {
		logs.PrintlnSuccess("sign ok .....................")
		return true
	}
	logs.PrintlnWarning("verify sign err ", pAddress.String(), address)
	return false
}

func FetchPubKeyBySignMsg(sign, msg string) (*ecdsa.PublicKey, error) {
	signBytes, err := hexutil.Decode(sign)
	if err != nil {
		return nil, err
	}

	if signBytes[64] != 27 && signBytes[64] != 28 {
		return nil, errors.New("invalid sign..")
	}
	signBytes[64] -= 27
	return cryptoEth.SigToPub(global.EthSignHash(msg), signBytes)
}

func GetFromStore(dir, name string) (*PrivateKey, error) {
	ks, err := keystore.NewFSKeystore(dir)
	if err != nil {
		return nil, err
	}
	pri, err := ks.Get(name)
	if err != nil {
		return nil, err
	}
	return NewPrivateKeyByLibP2pPri(pri)
}

func PubKeyToAddress(pubKey *ecdsa.PublicKey) common.Address {
	return cryptoEth.PubkeyToAddress(*pubKey)
}
