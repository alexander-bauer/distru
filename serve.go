package main

import (
	"bufio"
	"log"
	"net"
)

const (
	GETJSON = "distru json\r\n"  //Requests a json-encoded current index.
	NEWSITE = "distru index\r\n" //Prefaces a request to index a new site.
	SHARE   = "distru share\r\n" //Wraps Idx.MergeRemote()
	SAVE    = "distru save\r\n"  //Saves the current configuration and index
)

//The root dir should actually be a search page, which serves up a page to enter a search query, which is then turned into a search results page

//Serve is the primary function of distru. It listens on the tcp port 9049 for incoming connections, then passes them directly to handleConn.
func Serve(conf *config) {
	log.Println("Distru version", Version)
	log.Println("Configuration status:\n\tGenerated in:\t", conf.Version, "\n\tAutoIndexing:\t", len(conf.AutoIndex), "\n\tResouces:\t", len(conf.Resources), "\n\tSites indexed:\t", len(conf.Idx.Sites))

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
	defer conn.Close()
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
	case GETJSON:
		{
			//Then serve a json encoded index.
			_, err := w.WriteString(conf.Idx.JSON())
			if err != nil {
				log.Println(prefix, "error serving json:", err)
				return
			} //close if
			//and flush it to the connection.
			err = w.Flush()
			if err != nil {
				log.Println(prefix, "error serving json:", err)
				return
			} //close if
			conn.Close()
			log.Println(prefix, "served json")
		} //close case

	case NEWSITE:
		{
			siteRequest, err := r.ReadBytes('\n')
			if err != nil {
				log.Println(prefix, err)
				return
			}
			site := string(siteRequest[:len(siteRequest)-2])
			log.Println(prefix, "command to index:", site)
			conn.Close()
			conf.Idx.Queue <- site
		} //close case
	case SHARE:
		{
			shareRequest, err := r.ReadBytes('\n')
			if err != nil {
				log.Println(prefix, err)
				conn.Close()
			}
			conn.Close()
			remote := string(shareRequest[:len(shareRequest)-2])
			log.Println(prefix, "merging index from:", remote)
			err = conf.Idx.MergeRemote(remote, true, 0)
			if err != nil {
				log.Println(prefix, err)
			} else {
				log.Println(prefix, "merged from:", remote)
			}
		}
	case SAVE:
		{
			conn.Close()
			err := conf.save(ConfPath)
			if err != nil {
				log.Println(prefix, "error saving to:", ConfPath)
			}
			log.Println(prefix, "saved to:", ConfPath)
		}
	default:
		{
			//Display the request
			log.Println(prefix, "invalid request: \""+req+"\"")
			conn.Close()
		} //close default case
	} //close switch

}
