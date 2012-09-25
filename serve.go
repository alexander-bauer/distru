package main

import (
    "fmt"
    "net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hi there. This webpage is a response by <i>distru</i>, which is being run on the machine who's port you're looking at.\nGo is pretty cool.")
}

func serve() {
    http.HandleFunc("/", handler)
    http.ListenAndServe(":8080", nil)
}
