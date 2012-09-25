package main

import (
	"fmt"
	"net/http"
)

func handleReadable(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<html><body>Hi there. This webpage is a response by <i>Distru</i>, which is being run on the machine who's port you're looking at.\nGo is pretty cool.</body></html>")
}

func handleBinary(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<html><body>This page will be binary encoded for use by other instances of <i>Distru</i>.</body></html>")
}

func Serve() {
	http.HandleFunc("/", handleReadable)
	http.HandleFunc("/bin", handleBinary)
	http.ListenAndServe(":8080", nil)
}
