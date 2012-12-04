package main

import (
	"os"
)

const (
	Version = "0.10.8"
)

var (
	ConfPath = "/etc/distru.conf"
	Conf     *config
)

func main() {
	if len(os.Args) >= 2 {
		//Load from the first argument if it's
		//supplied.
		ConfPath = os.Args[1]
	}
	Conf = GetConfig(ConfPath)
	Serve(Conf)
}
