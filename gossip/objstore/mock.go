package objstore

import (
	"goworkshop/gossip/cas/swarmcas"
)

func NewMock() ObjectStore {
	return New(&Config{
		BackendStorage: swarmcas.NewMock(),
	})
}
