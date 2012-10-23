package main

import (
	"bufio"
	"encoding/gob"
	"log"
	"net"
)

//The root dir should actually be a search page, which serves up a page to enter a search query, which is then turned into a search results page

//Serve is the primary function of distru. It listens on the tcp port 9049 for incoming connections, then passes them directly to handleConn.
func Serve() {
	ln, err := net.Listen("tcp", ":9049")
	if err != nil {
		log.Fatal("Could not start server: ", err)
	}
	log.Println("Started server.")

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal("Server error.")
		}
		go handleConn(conn)
	}
}

//handleConn is the internal server function for distru. When it recieves a connection, it logs the RemoteAddr of the connection, then serves a gob of the in-memory index (Idx) to it. It closes the connection immediately afterward.
func handleConn(conn net.Conn) {
	w := bufio.NewWriter(conn)
	log.Println("Serving bin index to: " + conn.RemoteAddr().String())

	//Serve a gob to the new connection immediately.
	encoder := gob.NewEncoder(w)
	encoder.Encode(Idx)
}

//RecvIndex tries to recieve an index gob from a distru server (on tcp port 9049) running on the given url. It returns an empty index if it fails to do so.
func RecvIndex(url string) *Index {
	//Create the connection, from which the target server should immediately try to serve an index.
	conn, err := net.Dial("tcp", url+":9049")
	if err != nil {
		log.Println("No response from: " + url)
		return &Index{}
	}

	//When we're ready, create a reader, so we can retrieve the data from the connection.
	r := bufio.NewReader(conn)

	//Finally, try to use the gob decoder to form an index from the gob.
	decoder := gob.NewDecoder(r)
	index := &Index{}
	decoder.Decode(index)

	return index
}
