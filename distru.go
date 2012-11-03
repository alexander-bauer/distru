package main

const (
	Version = "0.3.3"
)

var Idx = &Index{Sites: make(map[string]*site)}

func main() {
	Serve()
}
