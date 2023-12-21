package main

import (
	"fmt"

	"source.local/common/pkg/env"
)

func main() {
	var envs struct {
		TestInt    int     `env:"TEST_INT"`
		TestString string  `env:"TEST_STRING"`
		TestBool   bool    `env:"TEST_BOOL"`
		TestFloat  float64 `env:"TEST_FLOAT"`
	}
	if err := env.Parse(&envs); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(envs)
}
