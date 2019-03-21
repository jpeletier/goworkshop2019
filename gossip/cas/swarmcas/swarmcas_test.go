package swarmcas_test

import (
	"bytes"
	"goworkshop/gossip/cas/swarmcas"
	"testing"
)

func TestGetPut(t *testing.T) {

	cas := swarmcas.NewMock()

	somedata := []byte("this is some data")

	key, err := cas.Put(somedata)
	if err != nil {
		t.Fatal(err)
	}

	expectedKey := "89512e57525313b220b73e61399e47dce11d5f704d1686cba77955aba2b5451a"

	if key != expectedKey {
		t.Fatalf("Expected key to be %s, got %s", expectedKey, key)
	}

	retrievedData, err := cas.Get(key)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(somedata, retrievedData) {
		t.Fatalf("Expected to retrieve %s, got %s", string(somedata), string(retrievedData))
	}

	// try to retrieve a key that does not exist:

	_, err = cas.Get("0000000000000000000000000000000000000000000000000000000000000000")

	if err == nil {
		t.Fatal("Expected to have an error since the key does not exist")
	}

}

// TODO: test with real Swarm
