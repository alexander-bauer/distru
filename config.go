package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

var (
	defaultConfig = &config{
		Version:    Version,
		IndexDelay: 60,
		IndexFile:  "/var/distru.index",
		WebDir:     "ui",
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
	IndexDelay int      //The number of minutes between Update() checks
	IndexFile  string   //The file to save the index to
	WebDir     string   //Directory containing stylesheets and webpages (including /)
	AutoIndex  []string //A list of sites to index on startup
	Resources  []string //A list of sites from which to request trusted indexes
	ResTimeout int      //Number of seconds to wait for response from Resources
	Idx        *Index   `json:"-"` //The local index
}

func (conf *config) save(filename string) error {
	b, err := json.MarshalIndent(conf, "", "\t")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, b, 0664)
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
	conf.Idx, err = LoadIndex(conf.IndexFile)
	if err != nil {
		conf.Idx = &Index{
			Sites: make(map[string]*site, 0),
			Cache: make([]*page, 0),
		}
		log.Println("Failed to load index from", conf.IndexFile+":", err)
	} else {
		log.Println("Loaded index from:", conf.IndexFile)
	}
	return conf
}
