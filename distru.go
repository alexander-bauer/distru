package main

const (
	Version = "0.2"
)

var Idx = &Index{Sites: make(map[string]site)}

func main() {
	Serve()
}
