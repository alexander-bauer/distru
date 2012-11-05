package main

const (
	Version  = "0.5.2"
	ConfPath = "/etc/distru.conf"
)

var Conf *config

func main() {
	Conf := GetConfig(ConfPath)
	Serve(Conf)
}
