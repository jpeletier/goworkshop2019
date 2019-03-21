package kv_test

import (
	"bytes"
	"goworkshop/gossip/kv"
	"goworkshop/gossip/kv/account"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/swarm/api/client"

	"github.com/ethereum/go-ethereum/crypto"

	"github.com/ethereum/go-ethereum/swarm/api"
	"github.com/ethereum/go-ethereum/swarm/storage"
	"github.com/ethereum/go-ethereum/swarm/storage/feed"

	swarmhttp "github.com/ethereum/go-ethereum/swarm/api/http"
)

func testkv(t *testing.T, account account.Account, kvservice kv.KV) {
	testkey := "testkey"
	testvalue := []byte("some content")

	if err := kvservice.Put(testkey, testvalue); err != nil {
		t.Fatal(err)
	}

	retrievedValue, err := kvservice.Get(account.Addr(), testkey)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(testvalue, retrievedValue) {
		t.Fatalf("Expected to retrieve %s, got %s", string(testvalue), string(retrievedValue))
	}

	_, err = kvservice.Get(account.Addr(), "fake key")
	if err != kv.ErrKeyNotFound {
		t.Fatalf("Expected Get to return ErrKeyNotFound when trying to look up a fake key")
	}
}

func TestKV(t *testing.T) {
	testServer, err := NewTestSwarmServer()
	if err != nil {
		t.Fatal(err)
	}
	defer testServer.Close()

	client := client.NewClient(testServer.URL)
	account := account.New("testpassphrase")

	kvservice := kv.New(&kv.Config{
		SwarmClient: client,
		Account:     account,
	})

	testkv(t, account, kvservice)

}

func TestKVMock(t *testing.T) {
	account := account.New("testpassphrase")

	kvservice := kv.NewMock(account)

	testkv(t, account, kvservice)

}

/*













 */

// swarm test utils -- disregard --

type TestServer interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
}

func serverFunc(api *api.API) TestServer {
	return swarmhttp.NewServer(api, "")
}

type TestSwarmServer struct {
	*httptest.Server
	Hasher    storage.SwarmHash
	FileStore *storage.FileStore
	dir       string
	cleanup   func()
}

//CreateDirectory ... Creates a directory and checks whether the action was successful
func CreateDirectory(path string) {

	_ = os.MkdirAll(path, 0750|os.ModeDir)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic(err)
	}

}

func NewTestSwarmServer() (*TestSwarmServer, error) {
	path, err := ioutil.TempDir("", "goworkshop")
	if err != nil {
		return nil, err
	}
	dir := filepath.Join(path, "main")
	CreateDirectory(dir)
	storeparams := storage.NewDefaultLocalStoreParams()
	storeparams.DbCapacity = 5000000
	storeparams.CacheCapacity = 5000
	storeparams.Init(dir)
	localStore, err := storage.NewLocalStore(storeparams, nil)
	if err != nil {
		os.RemoveAll(dir)
		return nil, err
	}
	fileStore := storage.NewFileStore(localStore, storage.NewFileStoreParams())

	// mutable resources test setup
	resourceDir := filepath.Join(path, "resources")
	CreateDirectory(resourceDir)

	rhparams := &feed.HandlerParams{}
	rh, err := feed.NewTestHandler(resourceDir, rhparams)
	if err != nil {
		return nil, err
	}
	pk, err := crypto.HexToECDSA("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	if err != nil {
		return nil, err
	}

	a := api.NewAPI(fileStore, nil, rh.Handler, pk)
	srv := httptest.NewServer(serverFunc(a))
	tss := &TestSwarmServer{
		Server:    srv,
		FileStore: fileStore,
		dir:       dir,
		Hasher:    storage.MakeHashFunc(storage.DefaultHash)(),
		cleanup: func() {
			srv.Close()
			rh.Close()
			os.RemoveAll(dir)
			os.RemoveAll(resourceDir)
		},
	}
	feed.TimestampProvider = tss
	return tss, err
}

func (t *TestSwarmServer) Close() {
	t.cleanup()
}

func (t *TestSwarmServer) Now() feed.Timestamp {
	return feed.Timestamp{Time: uint64(time.Now().Unix())}
}
