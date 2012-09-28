package main

import (
	"encoding/gob"
	"net"
	"bufio"
)

func fetchIndex(url string) *Index {
	conn, err := net.Dial("tcp", url + ":9049")
	if err != nil {
		print("Couldn't connect.")
		return &Index{}
	}
	
	r := bufio.NewReader(conn)
	
	decoder := gob.NewDecoder(r)
	index := &Index{}
	decoder.Decode(index)
	
	return index
}
