package kv

import (
	"errors"
	"goworkshop/gossip/kv/account"
	"io/ioutil"

	"github.com/ethereum/go-ethereum/crypto"

	"github.com/ethereum/go-ethereum/swarm/api/client"
	"github.com/ethereum/go-ethereum/swarm/storage/feed"

	"github.com/ethereum/go-ethereum/common"
)

type KV interface {
	Put(key string, value []byte) (err error)
	Get(addr common.Address, key string) (value []byte, err error)
}

type Config struct {
	SwarmClient *client.Client
	Account     account.Account
}

type kv struct {
	Config
}

var ErrKeyNotFound error = errors.New("Key not found")

func New(config *Config) KV {
	return &kv{
		Config: *config,
	}
}

func (k *kv) Put(key string, value []byte) (err error) {
	query := new(feed.Query)
	query.User = common.Address(k.Account.Addr())
	query.Topic = feed.Topic(crypto.Keccak256Hash([]byte(key)))
	request, err := k.SwarmClient.GetFeedRequest(query, "")
	if err != nil {
		return err
	}
	request.SetData(value)
	if err := request.Sign(k.Account.Signer()); err != nil {
		return err
	}

	return k.SwarmClient.UpdateFeed(request)
}

func (k *kv) Get(addr common.Address, key string) (value []byte, err error) {
	query := new(feed.Query)
	query.User = addr
	query.Topic = feed.Topic(crypto.Keccak256Hash([]byte(key)))
	reader, err := k.SwarmClient.QueryFeed(query, "")
	if err != nil {
		if err == client.ErrNoFeedUpdatesFound {
			return nil, ErrKeyNotFound
		}
		return nil, err
	}

	return ioutil.ReadAll(reader)
}
