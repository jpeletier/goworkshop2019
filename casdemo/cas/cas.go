package cas

type Cas interface {
	Put(data []byte) (key string, err error)
	Get(key string) (data []byte, err error)
}
