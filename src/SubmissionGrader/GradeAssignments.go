package main

import (
	"SubmissionGrader/internal"
	"fmt"
	"sync"
)

func main() {

	var wg sync.WaitGroup
	var sharedSignal internal.GoSafeVar[bool]

	fmt.Printf("Setting up os signal watcher.")
	go internal.WatchForSignal(&sharedSignal)

	fmt.Printf("Starting Grading assignment.")
	go internal.GradeAssignment(&sharedSignal, &wg)

	wg.Wait()

	fmt.Println("Grading Completed.")

	// done: start with environment variables
	// done: switch statement
	// TODO: one language at a time
	// TODO: java>c++>duplicate
	// done: sigterm handling on grader level
	// TODO: Test all functions (going forward do this before writing on when you abstract it out of a function)
}
