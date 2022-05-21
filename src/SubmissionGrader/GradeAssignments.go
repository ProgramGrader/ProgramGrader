package main

import (
	"SubmissionGrader/internal"
	"fmt"
	"sync"
)

func main() {

	var wg sync.WaitGroup
	var sharedSignal internal.GoSafeVar[bool]

	go internal.WatchForSignal(&sharedSignal)

	go internal.GradeAssignment(&sharedSignal, &wg)

	wg.Wait()

	fmt.Println("Exiting")

	// done: start with environment variables
	// done: switch statement
	// TODO: one language at a time
	// TODO: java>c++>duplicate
	// TODO: sigterm handling on grader level
}
