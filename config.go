package main

import (
	"encoding/json"
	"io/ioutil"
)

type config struct {
	Version    string   //The Distru version that generated this config
	IndexPath  string   //The path to which the Index is to be saved and loaded
	Indexers   int      //The number of indexer processes that should be run
	AutoIndex  []string //A list of sites to index on startup
	Resources  []string //A list of sites from which to request trusted indexes
	ResTimeout int      //Number of seconds to wait for response from Resources
	Idx        *Index   `json:"-"` //The local index, not saved with config
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
			IndexPath:  "/var/tmp/distru.index",
			Indexers:   1,
			AutoIndex:  make([]string, 0),
			Resources:  make([]string, 0),
			ResTimeout: 8,
		}
		conf.save(filename)
	}
	conf.Idx, _ = LoadIndex(conf.IndexPath)
	return conf
}
