package timeline

import (
	"goworkshop/gossip/kv"
	"goworkshop/gossip/kv/account"
	"goworkshop/gossip/objstore"
)

func NewMock(account account.Account) Timeline {
	return New(&Config{
		KVService: kv.NewMock(account),
		ObjStore:  objstore.NewMock(),
		Account:   account,
	})
}
