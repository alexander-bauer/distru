package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

var (
	defaultConfig = &config{
		Version:    Version,
		Indexers:   1,
		WebDir:     "ui/",
		AutoIndex:  make([]string, 0),
		Resources:  make([]string, 0),
		ResTimeout: 8,
		Idx: &Index{
			Sites: make(map[string]*site, 0),
			Cache: make([]*page, 0),
		},
	}
)

type config struct {
	Version    string   //The Distru version that generated this config
	Indexers   int      //The number of indexer processes that should be run
	WebDir     string   //Directory containing stylesheets and webpages (including /)
	AutoIndex  []string //A list of sites to index on startup
	Resources  []string //A list of sites from which to request trusted indexes
	ResTimeout int      //Number of seconds to wait for response from Resources
	Idx        *Index   `json:",omitempty"` //The local index
}

func (conf *config) save(filename string) error {
	b, err := json.MarshalIndent(conf, "", "\t")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, b, 0660)
}

func loadConf(filename string) (conf *config, err error) {
	//If we get an error, return the default.
	conf = &*defaultConfig

	var b []byte
	b, err = ioutil.ReadFile(filename)
	if err != nil {
		return
	}
	err = json.Unmarshal(b, &conf)
	if err != nil {
		return
	}
	return
}

func GetConfig(filename string) *config {
	conf, err := loadConf(filename)
	if err != nil {
		//If we got an error, then we were also
		//given the default, so save it.
		log.Println("Failed to load config:", err)
		log.Println("Saving default config to:", filename)
		err = conf.save(filename)
		if err != nil {
			log.Println("Error saving config:", err)
			log.Println("Using default config anyway")
		}
	}
	return conf
}
