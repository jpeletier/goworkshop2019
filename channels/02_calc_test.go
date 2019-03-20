package channels_test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func slowOperation(a, b int, result chan<- int) {
	randomDelay := time.Duration(rand.Intn(3)+1) * time.Second
	time.Sleep(randomDelay)
	result <- a * b
}

func TestChannelCalc(t *testing.T) {
	result := make(chan int)

	numTasks := 10

	for n := 0; n < numTasks; n++ {
		go slowOperation(n, n+1, result)
	}

	for n := 0; n < numTasks; n++ {
		r := <-result
		fmt.Printf("Result received: %d\n", r)
	}

}
