package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
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
	numResults := 2
	log.Println("<-" + r.RemoteAddr + "> searching \"" + searchTerm + "\"")

	//load external files
	css, err := ioutil.ReadFile("ui/search.css")
	if err != nil {
		panic(err)
	}

	//add the page
	w.Write([]byte("<html><head><title>Distru :: Searching " + searchTerm + "</title><style type=\"text/css\">"))
	w.Write(css)
	w.Write([]byte("</style></head><body><div class=\"searchterm\">" + strconv.Itoa(numResults) + " results for <strong>" + searchTerm + "</strong></div>"))

	//TODO: SEARCH HERE.
	//THIS IS JUST AN EXAMPLE..
	w.Write([]byte("<div class=\"results\">test</div>" + "</body></html>"))
}

func frontpageHandler(w http.ResponseWriter, r *http.Request) {
	//load external files
	css, err := ioutil.ReadFile("ui/index.css")
	if err != nil {
		panic(err)
	}
	javascript, err := ioutil.ReadFile("ui/distru.js")
	if err != nil {
		panic(err)
	}

	//add the page
	w.Write([]byte("<html><head><title>Distru :: Search Freely</title><style type=\"text/css\">"))
	w.Write(css)
	w.Write([]byte("</style><script type=\"text/javascript\">"))
	w.Write(javascript)
	w.Write([]byte("</script></head><body><div class = \"version\">"))
	w.Write([]byte(Version))
	w.Write([]byte("<div class = \"name\">Distru</div><input type=\"text\" onkeydown=\"searchThis();\" onkeypress=\"isEnter(event);\" id=\"search\" class=\"search\" placeholder=\"Search freely\"/></body></html>"))
}
