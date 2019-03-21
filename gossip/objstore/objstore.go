package objstore

import (
	"encoding/json"
	"goworkshop/gossip/cas"
)

type ObjectStore interface {
	Put(obj interface{}) (key string, err error)
	Get(key string, obj interface{}) (err error)
}

type Config struct {
	BackendStorage cas.Cas
}

type objectStore struct {
	Config
}

func New(config *Config) ObjectStore {
	return &objectStore{
		Config: *config,
	}
}

func (os *objectStore) Put(obj interface{}) (key string, err error) {
	objBytes, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}
	return os.BackendStorage.Put(objBytes)
}

func (os *objectStore) Get(key string, obj interface{}) (err error) {
	objBytes, err := os.BackendStorage.Get(key)
	if err != nil {
		return err
	}

	return json.Unmarshal(objBytes, obj)
}
