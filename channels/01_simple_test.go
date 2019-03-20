package channels_test

import (
	"fmt"
	"testing"
	"time"
)

func verySlowFunction() {
	fmt.Println("very slow function launched")
	time.Sleep(5 * time.Second)
	fmt.Println("very slow function finished")

}

func waitWithMessage(done chan struct{}) {
	for {
		select {
		case <-done:
			fmt.Println("Finished at last!!")
			return
		default:
			fmt.Println("Let's wait a bit more")
			time.Sleep(1 * time.Second)
		}
	}
}

func TestChannelHelloWorld(t *testing.T) {
	done := make(chan struct{})

	go func() {
		verySlowFunction()
		close(done)
	}()

	waitWithMessage(done)

	fmt.Println("End of program")
}
