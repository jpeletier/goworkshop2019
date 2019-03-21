package timeline

import (
	"goworkshop/gossip/kv"
	"goworkshop/gossip/kv/account"
	"goworkshop/gossip/objstore"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

type Comment struct {
	Text      string `json:"content"`
	Timestamp int64  `json:"timestamp"`
	Previous  string `json:"previous"`
}

type Timeline interface {
	Post(text string) error
	Dump(addr common.Address) (comments <-chan *Comment)
}

type Config struct {
	KVService kv.KV
	ObjStore  objstore.ObjectStore
	Account   account.Account
}

type timeline struct {
	Config
}

const lastPostKey = "LASTPOST"

func New(config *Config) Timeline {
	return &timeline{
		Config: *config,
	}
}

func (t *timeline) Post(text string) error {
	// get current post if it exists
	keyBytes, err := t.KVService.Get(t.Account.Addr(), lastPostKey)
	if err != nil && err != kv.ErrKeyNotFound {
		return err
	}

	comment := &Comment{
		Previous:  string(keyBytes),
		Timestamp: time.Now().Unix(),
		Text:      text,
	}

	key, err := t.ObjStore.Put(comment)
	if err != nil {
		return err
	}
	if err := t.KVService.Put(lastPostKey, []byte(key)); err != nil {
		return err
	}
	return nil
}

func (t *timeline) Dump(addr common.Address) <-chan *Comment {
	comments := make(chan *Comment)
	go func() {
		defer close(comments)
		keyBytes, err := t.KVService.Get(addr, lastPostKey)
		if err != nil && err != kv.ErrKeyNotFound {
			return
		}
		for len(keyBytes) > 0 {
			var comment Comment
			if err := t.ObjStore.Get(string(keyBytes), &comment); err != nil {
				return
			}
			comments <- &comment
			keyBytes = []byte(comment.Previous)
		}
	}()
	return comments

}
