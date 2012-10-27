package main

import "os"

var Idx = &Index{Sites: make(map[string]site)}

func main() {
	Serve()
	os.Exit(1)
}
