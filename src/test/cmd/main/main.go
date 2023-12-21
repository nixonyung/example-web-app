package main

import (
	"errors"
	"sync"

	"source.local/common/pkg/servicebase"
	"source.local/test/internal/testcases"
)

func main() {
	testcases := [](func() error){
		testcases.TestUsers,
	}

	// (ref.) [Golang concurrency: how to append to the same slice from different goroutines](Golang concurrency: how to append to the same slice from different goroutines)
	testResults := make([]error, len(testcases))
	var wg sync.WaitGroup
	for i, test := range testcases {
		wg.Add(1)
		go func(i int, test func() error) {
			testResults[i] = test()
			wg.Done()
		}(i, test)
	}
	wg.Wait()
	servicebase.HandleErr(errors.Join(testResults...))
}
