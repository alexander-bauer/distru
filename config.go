package main

import (
	"encoding/json"
	"io/ioutil"
)

type config struct {
	Version   string   //The Distru version that generated this config
	Indexers  int      //The number of indexer processes that should be run
	AutoIndex []string //A list of sites to index on startup
	Resources []string //A list of sites from which to request trusted indexes
	Idx       *Index   `json:",omitempty"` //The local index
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
			Version: Version,
			Idx: &Index{
				Sites: make(map[string]*site),
			},
		}
	}
	return conf
}
