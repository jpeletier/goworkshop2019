package marshalling_test

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
	"time"
)

type SomeStruct struct {
	A string
	B int
	C bool
}

type Person struct {
	Name      string    `json:"name"`
	BirthDate time.Time `json:"age"`
	Cool      bool      `json:"cool"`
}

func TestMarshalling1(t *testing.T) {
	data := &SomeStruct{
		A: "Hola",
		B: 1979,
		C: true,
	}
	b, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(string(b))
}

func TestMarshalling2(t *testing.T) {
	person := &Person{
		Name:      "Javi",
		BirthDate: time.Date(1979, time.April, 4, 9, 30, 0, 0, time.Local),
		Cool:      true,
	}
	personBytes, err := json.Marshal(person)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(string(personBytes))

	var person2 Person

	if err := json.Unmarshal(personBytes, &person2); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(person, &person2) {
		t.Fatal(err)
	}
}

type Timestamp time.Time

type Person2 struct {
	Name      string    `json:"name"`
	BirthDate Timestamp `json:"age"`
	Cool      bool      `json:"cool"`
}

func (t Timestamp) MarshalJSON() ([]byte, error) {
	tm := time.Time(t).Unix()
	return json.Marshal(tm)
}

func (t *Timestamp) UnmarshalJSON(data []byte) error {
	var tm int64
	if err := json.Unmarshal(data, &tm); err != nil {
		return err
	}
	*t = Timestamp(time.Unix(tm, 0))
	return nil
}

func TestCustomMarshalling(t *testing.T) {
	person := &Person2{
		Name:      "Javi",
		BirthDate: Timestamp(time.Date(1979, time.April, 4, 9, 30, 0, 0, time.Local)),
		Cool:      true,
	}
	personBytes, err := json.Marshal(person)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(string(personBytes))

	var person2 Person2

	if err := json.Unmarshal(personBytes, &person2); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(person, &person2) {
		t.Fatal(err)
	}
}
