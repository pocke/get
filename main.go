package main

import (
	"fmt"
	"os"
)

func main() {
	if err := Main(os.Args); err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
}

func Main(args []string) error {
	if len(args) < 3 {
		return fmt.Errorf("Too few arguments. Please `get TYPE ADDR`")
	}
	t := args[1]
	addr := args[2]

	fn, ok := Getters[t]
	if !ok {
		return fmt.Errorf("Type %s doesn't exist", t)
	}

	return fn(addr)
}

var Getters = map[string]func(addr string) error{
	"ghq": func(addr string) error {
		return nil
	},
	"go": func(addr string) error {
		return nil
	},
}
