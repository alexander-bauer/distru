package main

import (
	"fmt"
	"net/http"
)

//the root dir should actually be a search page, which serves up a page to enter a search query, which is then turned into a search results page

func handleReadable(w http.ResponseWriter, r *http.Request) {
	s := RepIndex(NewIndex())
	fmt.Fprintf(w, s)
}

func handleBinary(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<html><body>This page will be binary encoded for use by other instances of <i>Distru</i>.</body></html>")
}

func Serve() {
	http.HandleFunc("/index/text", handleReadable)
	http.HandleFunc("/index/bin", handleBinary)
	http.ListenAndServe(":9049", nil)
}
