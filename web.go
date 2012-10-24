package main

import (
	"fmt"
	"log"
	"net/http"
)

func ServeWeb() {
	http.HandleFunc("/", searchHandler)
	log.Println("Starting webserver.")
	http.ListenAndServe(":9048", nil)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	s := "placeholder"

	fmt.Fprintf(w, s)
}
