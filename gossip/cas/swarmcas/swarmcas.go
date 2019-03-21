package swarmcas

import (
	"bytes"
	"goworkshop/gossip/cas"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

func init() {
	http.DefaultClient.Timeout = 10 * time.Second // hack to have Swarm Client time out if it can't find a key
}

type SwarmClientRaw interface {
	UploadRaw(io.Reader, int64, bool) (string, error)
	DownloadRaw(key string) (io.ReadCloser, bool, error)
}

type Config struct {
	SwarmClient SwarmClientRaw
}

type swarmCas struct {
	Config
}

func New(config Config) cas.Cas {
	return &swarmCas{
		Config: config,
	}
}

func (s *swarmCas) Put(data []byte) (key string, err error) {
	r := bytes.NewReader(data)
	return s.SwarmClient.UploadRaw(r, int64(r.Len()), false)
}

func (s *swarmCas) Get(key string) (data []byte, err error) {
	r, _, err := s.SwarmClient.DownloadRaw(key)
	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(r)
}
