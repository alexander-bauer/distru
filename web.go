package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func ServeWeb() {
	http.HandleFunc("/search/", searchHandler)
	http.HandleFunc("/", frontpageHandler)
	log.Println("Started webserver on port 9048.")
	http.ListenAndServe(":9048", nil)

}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	//get the search term and save it as searchTerm
	searchTerm := r.URL.Path[len("/search/"):]
	//get the number of results for the searchTerm
	numResults := 0
	log.Println("<-" + r.RemoteAddr + "> searching \"" + searchTerm + "\"")
	
	//add the page
	fmt.Fprintf(w, "<html><head><title>Distru :: Searching \"%s\"</title><style type=\"text/css\">", searchTerm)
	file, err := ioutil.ReadFile("ui/search.css")
	if err != nil {
		panic(err)
	}
	fmt.Fprint(w, string(file))
	fmt.Fprintf(w, "</style></head><body><div class=\"searchterm\">%d results for <strong>%s</strong></div>", numResults, searchTerm)

	//TODO: SEARCH HERE.
	//this is a temporary example of what searches will look like
	fmt.Fprint(w, "<div class=\"results\">test</div>")
	fmt.Fprint(w, "<div class=\"results\">test2</div>")
	
	//close page
	fmt.Fprint(w, "</body></html>")
}

func frontpageHandler(w http.ResponseWriter, r *http.Request) {
	//add the page
	fmt.Fprint(w, "<html><head><title>Distru :: Search Freely</title><style type=\"text/css\">")
	file, err := ioutil.ReadFile("ui/index.css")
	if err != nil {
		panic(err)
	}
	fmt.Fprint(w, string(file))
	fmt.Fprint(w, "</style><script type=\"text/javascript\">")
	file, err = ioutil.ReadFile("ui/distru.js")
	if err != nil {
		panic(err)
	}
	fmt.Fprint(w, string(file))
	fmt.Fprintf(w, "</script></head><body><div class = \"version\">Version %s</div><div class = \"name\">Distru</div><input type=\"text\" onkeydown=\"searchThis();\" onkeypress=\"isEnter(event);\" id=\"search\" class=\"search\" placeholder=\"Search freely\"/></body></html>", Version)
}
