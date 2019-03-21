package account

import (
	"crypto/ecdsa"
	"encoding/json"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/swarm/storage/feed"
)

type Signer interface {
	Sign(common.Hash) (feed.Signature, error)
	Address() common.Address
}

type Account interface {
	Addr() common.Address
	Signer() Signer
}

type accountWithPrivateKey struct {
	addr common.Address
	pk   *ecdsa.PrivateKey
}

func New(passphrase string) Account {
	d := crypto.Keccak256([]byte(passphrase))
	pk, _ := crypto.ToECDSA(d)
	return &accountWithPrivateKey{
		addr: common.Address(crypto.PubkeyToAddress(pk.PublicKey)),
		pk:   pk,
	}
}

func (a *accountWithPrivateKey) Addr() common.Address {
	return a.addr
}

func (a *accountWithPrivateKey) Signer() Signer {
	return feed.NewGenericSigner(a.pk)
}

func (a *accountWithPrivateKey) MarshalJSON() ([]byte, error) {
	st := a.Addr().Hex()
	return json.Marshal(st)
}
