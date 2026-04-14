package wallet

import (
	"crypto/sha256"
	"fmt"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil/base58"
	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/tyler-smith/go-bip39"
	"golang.org/x/crypto/sha3"
)

const (
	tronPurpose  = uint32(44)
	tronCoinType = uint32(195)
	tronAccount  = uint32(0)
	tronChange   = uint32(0)
)

func DeriveTRONAddress(mnemonic string, index uint32) (string, string, error) {
	if !bip39.IsMnemonicValid(mnemonic) {
		return "", "", fmt.Errorf("%w: invalid mnemonic", ErrInvalidArgument)
	}

	seed := bip39.NewSeed(mnemonic, "")
	key, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	if err != nil {
		return "", "", err
	}

	path := fmt.Sprintf("m/44'/195'/%d'/%d/%d", tronAccount, tronChange, index)
	key, err = derive(key, hdkeychain.HardenedKeyStart+tronPurpose)
	if err != nil {
		return "", "", err
	}
	key, err = derive(key, hdkeychain.HardenedKeyStart+tronCoinType)
	if err != nil {
		return "", "", err
	}
	key, err = derive(key, hdkeychain.HardenedKeyStart+tronAccount)
	if err != nil {
		return "", "", err
	}
	key, err = derive(key, tronChange)
	if err != nil {
		return "", "", err
	}
	key, err = derive(key, index)
	if err != nil {
		return "", "", err
	}

	pubKey, err := key.ECPubKey()
	if err != nil {
		return "", "", err
	}

	address := tronAddressFromPubKey(pubKey)
	return address, path, nil
}

func derive(key *hdkeychain.ExtendedKey, index uint32) (*hdkeychain.ExtendedKey, error) {
	child, err := key.Derive(index)
	if err == nil {
		return child, nil
	}

	return nil, err
}

func tronAddressFromPubKey(pubKey *btcec.PublicKey) string {
	uncompressed := pubKey.SerializeUncompressed()
	hash := sha3.NewLegacyKeccak256()
	hash.Write(uncompressed[1:])
	sum := hash.Sum(nil)

	payload := append([]byte{0x41}, sum[len(sum)-20:]...)

	first := sha256.Sum256(payload)
	second := sha256.Sum256(first[:])
	encoded := append(payload, second[:4]...)

	return base58.Encode(encoded)
}
