package main

import (
	"os"

	"github.com/beeceej/structural"
)

func main() {
	if err := structural.Generate(os.Stdout, "helloworld/definition.go"); err != nil {
		panic(err.Error())
	}
}
