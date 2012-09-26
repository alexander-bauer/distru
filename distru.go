package main

import "os"

var Idx = NewIndex()

func main() {
	Serve()
	os.Exit(0)
}
