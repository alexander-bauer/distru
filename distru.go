package main

import "os"

var Idx = NewIndex()

func main() {
	if os.Args[1] == "serve" {
		Serve()
	}
	s := fetch(os.Args[1])
	links := getLinks(s)

	internal := getInternalLinks(links, s)
	external := getExternalLinks(links)

	print("Internal Links:\n")
	for i := range internal {
		print(internal[i], "\n")
	}

	print("\n")

	print("External Links:\n")
	for i := range external {
		print(external[i], "\n")
	}

	os.Exit(0)
}
