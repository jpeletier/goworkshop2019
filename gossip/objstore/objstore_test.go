package objstore_test

import (
	"goworkshop/gossip/objstore"
	"reflect"
	"testing"
)

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
	Cool bool   `json:"cool"`
}

func TestObjStore(t *testing.T) {
	store := objstore.NewMock()

	testobj := &Person{
		Name: "Perry Mason",
		Age:  39,
		Cool: true,
	}

	expectedKey := "c831f0f1c406653f4d5a4b63f14405cc61d5f0070ee4cf4601032e6ca711cb7d"
	key, err := store.Put(testobj)
	if err != nil {
		t.Fatal(err)
	}
	if key != expectedKey {
		t.Fatalf("Expected key to be %s, got %s", expectedKey, key)
	}

	var retrieved Person

	if err := store.Get(key, &retrieved); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(testobj, &retrieved) {
		t.Fatalf("Expected retrieved object to be %v, got %v", testobj, &retrieved)
	}

	// attempt to retrieve a non-existing object:
	if err := store.Get("0000000000000000000000000000000000000000000000000000000000000000", &retrieved); err == nil {
		t.Fatalf("Expected to receive an error since key does not exist")
	}

}
