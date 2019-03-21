package channels_test

import (
	"fmt"
	"testing"
	"time"
)

type CalcJob struct {
	id, a, b, result int
}

func runCalculator(input <-chan *CalcJob, output chan<- *CalcJob) {
	fmt.Println("Calculator server started")
	for job := range input {
		fmt.Printf("Starting job #%d\n", job.id)
		time.Sleep(1 * time.Second)
		job.result = job.a * job.b
		output <- job
	}
	fmt.Println("Calculator server finished")

}

func runPrintResults(input chan<- *CalcJob, output <-chan *CalcJob) {
	fmt.Println("Result printer started")
	for job := range output {
		fmt.Printf("Received result %d from job #%d\n", job.result, job.id)
	}
	fmt.Println("Result printer finished")

}

func TestCalculatorServer(t *testing.T) {
	input := make(chan *CalcJob)
	output := make(chan *CalcJob)

	numJobs := 10

	go runCalculator(input, output)
	go runPrintResults(input, output)

	for n := 0; n < numJobs; n++ {
		fmt.Printf("Submitting job #%d...\n", n)
		job := &CalcJob{
			id: n,
			a:  n,
			b:  n + 7,
		}
		input <- job
	}

	close(input)

	time.Sleep(5 * time.Second)
	close(output)
}

func TestCalculatorServerInputBuffer(t *testing.T) {
	input := make(chan *CalcJob, 5)
	output := make(chan *CalcJob)

	numJobs := 10

	go runCalculator(input, output)
	go runPrintResults(input, output)

	for n := 0; n < numJobs; n++ {
		fmt.Printf("Submitting job #%d...\n", n)
		job := &CalcJob{
			id: n,
			a:  n,
			b:  n + 7,
		}
		input <- job
	}

	close(input)

	time.Sleep(10 * time.Second)
}
