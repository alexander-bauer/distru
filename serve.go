package main

import (
	"log"
	"net"
)

//The root dir should actually be a search page, which serves up a page to enter a search query, which is then turned into a search results page

//Serve is the primary function of distru. It listens on the tcp port 9049 for incoming connections, then passes them directly to handleConn.
func Serve(conf *config) {
	log.Println("Distru version", Version)
	log.Println("Configuration status:\n\tGenerated in:\t", conf.Version,
		"\n\tIndexers:\t", conf.Indexers,
		"\n\tAutoIndexing:\t", len(conf.AutoIndex),
		"\n\tResouces:\t", len(conf.Resources),
		"\n\tSites indexed:\t", len(conf.Idx.Sites))

	//Start the Index Maintainer for the index.
	MaintainIndex(conf.Idx, conf.Indexers)

	go func() {
		for i := range conf.AutoIndex {
			conf.Idx.Queue <- conf.AutoIndex[i]
		}
	}()

	ln, err := net.Listen("tcp", ":9049")
	if err != nil {
		log.Fatal("Could not start server:", err)
	}
	log.Println("Started server on port 9049.")

	//Start a new goroutine for the webserver.
	go ServeWeb()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal("Server error.")
		}
		go handleConn(conf, conn)
	}
}

//handleConn is the internal server function for distru. When it recieves a connection, it waits for an instruction such as "distru json". It responds, then closes the connection.
func handleConn(conf *config, conn net.Conn) {
	//Ensure that the connection closes.
	defer conn.Close()

	//Simplify the logging.
	prefix := "<-" + conn.RemoteAddr().String() + ">"

	log.Println(prefix, "new connection")
	err := conf.Idx.JSON(conn)
	if err != nil {
		log.Println(prefix, err)
	} else {
		log.Println(prefix, "served index")
	}
	conn.Close()

}
