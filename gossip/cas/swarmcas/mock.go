package swarmcas

import "goworkshop/gossip/cas"

func NewMock() cas.Cas {
	return New(Config{
		SwarmClient: NewSwarmClientMock(),
	})
}
