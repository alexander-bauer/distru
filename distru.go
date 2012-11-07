package main

const (
	Version  = "0.7.2"
	ConfPath = "/etc/distru.conf"
)

var Conf *config

func main() {
	Conf = GetConfig(ConfPath)
	Serve(Conf)
}
