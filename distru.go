package main

const (
	Version = "0.3.4"
)

var Idx = &Index{Sites: make(map[string]*site)}

func main() {
	Serve()
}
