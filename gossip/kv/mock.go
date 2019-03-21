package kv

import (
	"fmt"
	"goworkshop/gossip/kv/account"

	"github.com/ethereum/go-ethereum/common"
)

type mock struct {
	table       map[string][]byte
	accountAddr string
}

func NewMock(account account.Account) KV {
	return &mock{
		table:       make(map[string][]byte),
		accountAddr: account.Addr().Hex(),
	}
}

func (m *mock) Put(key string, value []byte) (err error) {
	key = fmt.Sprintf("%s-%s", m.accountAddr, key)
	m.table[key] = value
	return nil
}

func (m *mock) Get(addr common.Address, key string) (value []byte, err error) {
	key = fmt.Sprintf("%s-%s", addr.Hex(), key)
	value = m.table[key]
	if value == nil {
		return nil, ErrKeyNotFound
	}
	return value, nil
}
