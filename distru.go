package main

const (
	Version = "0.5.1"
)

var Idx = &Index{Sites: make(map[string]*site)}

func main() {
	Serve()
}
