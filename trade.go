package main

import (
	"encoding/gob"
	"net"
	"bufio"
	"io"
	"log"
)

func RecvIndex(url string) *Index {
	conn, err := net.Dial("tcp", url + ":9049")
	if err != nil {
		log.Println("No response from: " + url)
		return &Index{}
	}
	
	r := bufio.NewReader(conn)
	
	decoder := gob.NewDecoder(r)
	index := &Index{}
	decoder.Decode(index)
	
	return index
}

func ServIndex(w io.Writer, index *Index) {
	encoder := gob.NewEncoder(w)
	encoder.Encode(index)
}
