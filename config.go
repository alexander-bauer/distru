package main

import (
	"encoding/json"
	"io/ioutil"
)

const (
	//Default values for the config type
	DIndexers   = 1
	DResTimeout = 8
)

type config struct {
	Version    string   //The Distru version that generated this config
	Indexers   int      //The number of indexer processes that should be run
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

func loadConf(filename string) (*config, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var conf config
	err = json.Unmarshal(b, &conf)
	if err != nil {
		return nil, err
	}
	return &conf, nil
}

func GetConfig(filename string) *config {
	conf, err := loadConf(filename)
	if err != nil {
		conf = &config{
			Version:    Version,
			Indexers:   DIndexers,
			AutoIndex:  make([]string, 0),
			Resources:  make([]string, 0),
			ResTimeout: DResTimeout,
			Idx: &Index{
				Sites: make(map[string]*site),
				Cache: make(map[string]*page),
			},
		}
	}
	return conf
}
