package main

import (
	"sync"

	"source.local/test/internal/testcases"
)

func main() {
	var wg sync.WaitGroup
	for _, testcase := range [](func()){
		testcases.TestUsersCRUD,
	} {
		wg.Add(1)
		go func(testcase func()) {
			testcase()
			wg.Done()
		}(testcase)
	}
	wg.Wait()
}
