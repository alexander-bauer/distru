package main

import "os"

var Idx = &Index{Sites: make(map[string]site)}
var Queue = make(chan<- string)

func main() {
	Serve()
	os.Exit(1)
}
