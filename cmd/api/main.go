package main

import (
	"net/http"

	"github.com/beeceej/structural/helloworld"
)

func main() {
	mux := http.DefaultServeMux
	api := new(helloworld.HelloWorldAPI)
	api.Routes(mux)
	err := http.ListenAndServe(":8000", mux)

	if err != nil {
		panic(err.Error())
	}
}
