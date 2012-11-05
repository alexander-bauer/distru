package main

const (
	Version  = "0.6"
	ConfPath = "/etc/distru.conf"
)

var Conf *config

func main() {
	Conf = GetConfig(ConfPath)
	Serve(Conf)
}
