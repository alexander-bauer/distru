package main

import (
	"bufio"
	"encoding/gob"
	"log"
	"net"
)

const (
	GETGOB  = "distru gob\r\n"   //Requests a gob of the current index.
	GETJSON = "distru json\r\n"  //Requests a json-encoded current index.
	NEWSITE = "distru index\r\n" //Prefaces a request to index a new site.
)

//The root dir should actually be a search page, which serves up a page to enter a search query, which is then turned into a search results page

//Serve is the primary function of distru. It listens on the tcp port 9049 for incoming connections, then passes them directly to handleConn.
func Serve() {
	ln, err := net.Listen("tcp", ":9049")
	if err != nil {
		log.Fatal("Could not start server:", err)
	}
	log.Println("Started server on port 9049.")

	//Start a new goroutine for the webserver.
	go ServeWeb()

	//Start the Index Maintainer, and recieve the input channel for it.
	MaintainIndex(Idx, 1)

	//Put a new domain into the queue.
	Idx.Queue <- "example.com"

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
	//Save the connection detail for simplicity of logging.
	prefix := "<-" + conn.RemoteAddr().String() + ">"
	log.Println(prefix, "new connection")

	//Going to check the request here, so create a new reader and writer
	r := bufio.NewReader(conn)
	w := bufio.NewWriter(conn)
	//and then read until we get a '\n', which should be preceded by '\r'
	b, err := r.ReadBytes('\n')
	if err != nil {
		log.Println(prefix, err)
		conn.Close()
	}
	req := string(b)
	
	switch req {
		case GETGOB: {
			Idx.Gob(w)
			conn.Close()
			log.Println(prefix, "served gob")
		} //close case
		
		case GETJSON: {
			//Then serve a json encoded index.
			_, err := w.WriteString(Idx.JSON())
			if err != nil {
				log.Println(prefix, "error serving json:", err)
				conn.Close()
				return	
			} //close if
			//and flush it to the connection.
			err = w.Flush()
			if err != nil {
				log.Println(prefix, "error serving json:", err)
				conn.Close()
				return
			} //close if
			conn.Close()
			log.Println(prefix, "served json")
		} //close case
		
		case NEWSITE: {
			site, err := r.ReadBytes('\n')
			if err != nil {
				log.Println(prefix, err)
				conn.Close()
			}
			Idx.Queue <- string(site[:len(site)-2])
			conn.Close()
		} //close case
		
		default: {
			//Display the request
			log.Println(prefix, "invalid request: \""+req+"\"")
			conn.Close()
		} //close default case
	} //close switch
	
}

//RecvIndex tries to recieve an index gob from a distru server (on tcp port 9049) running on the given url. It returns an empty index if it fails to do so.
func RecvIndex(url string) *Index {
	//Create the connection, from which the target server should immediately try to serve an index.
	conn, err := net.Dial("tcp", url+":9049")
	if err != nil {
		log.Println("No response from: " + url)
		return &Index{}
	}

	prefix := "->" + conn.RemoteAddr().String() + ">"

	//When we're ready, create a reader, so we can retrieve the data from the connection, and a writer, so we can request it.
	r := bufio.NewReader(conn)
	w := bufio.NewWriter(conn)

	//Request a gob from the target server
	log.Println(prefix, "requesting gob")
	_, err = w.WriteString(GETGOB)
	if err != nil {
		log.Println(prefix, "connection problem:", err)
		conn.Close()
		return &Index{}
	}
	err = w.Flush()
	if err != nil {
		log.Println(prefix, "connection problem:", err)
		conn.Close()
		return &Index{}
	}

	//Finally, try to use the gob decoder to form an index from the gob.
	decoder := gob.NewDecoder(r)
	index := &Index{}
	decoder.Decode(index)

	return index
}
