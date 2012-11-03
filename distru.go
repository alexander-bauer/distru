package main

const (
	Version = "0.3.1"
)

var Idx = &Index{Sites: make(map[string]*site)}

func main() {
	Serve()
}
