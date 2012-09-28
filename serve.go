package main

import (
	"bufio"
	"log"
	"net"
)

//the root dir should actually be a search page, which serves up a page to enter a search query, which is then turned into a search results page

func handleBinary(conn net.Conn) {
	w := bufio.NewWriter(conn)
	log.Println("Serving bin index to: " + conn.RemoteAddr().String())
	ServIndex(w, Idx)
}

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
		go handleBinary(conn)
	}
}
