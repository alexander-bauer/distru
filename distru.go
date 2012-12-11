package main

import (
	"flag"
	"log"
)

const (
	Version = "0.12"
)

var (
	ConfPath = "/etc/distru.conf"
	Conf     *config
)

var (
	optDefLoc  = flag.Bool("confloc", false, "print the default config-load location and exit")
	optGenConf = flag.Bool("genconf", false, "create a default config file")
	optWebDir  = flag.String("webdir", "", "WebDir to create a config with")
)

func main() {
	flag.Parse()

	if *optDefLoc {
		println(ConfPath)
		return
	}

	if flag.NArg() >= 1 {
		//Load from the first argument if it's
		//supplied.
		ConfPath = flag.Arg(0)
	}

	if *optGenConf {
		Conf = defaultConfig
		if len(*optWebDir) != 0 {
			Conf.WebDir = *optWebDir
		}
		err := defaultConfig.save(ConfPath)
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	Conf = GetConfig(ConfPath)
	Serve(Conf)
}
