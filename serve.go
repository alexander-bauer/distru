package main

import (
	"net"
	"bufio"
	"os"
	"log"
)

//the root dir should actually be a search page, which serves up a page to enter a search query, which is then turned into a search results page

func handleBinary(conn net.Conn) {
	w := bufio.NewWriter(conn)

	BinIndex(w, Idx)
}

func Serve() {
	ln, err := net.Listen("tcp", ":9049")
	if err != nil {
		print("Could not start server.")
		os.Exit(1)
	}
	print("Started server.\n")
	
	for {
		conn, err := ln.Accept()
		if err != nil {
			print("Server error.")
			os.Exit(1)
		}
		go handleBinary(conn)
	}
}
