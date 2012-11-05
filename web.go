package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
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
	log.Println("<-" + r.RemoteAddr + "> searching \"" + searchTerm + "\"")

	//Perform the search.
	num, results := Idx.Search(strings.Split(searchTerm, " "), 10)

	log.Println("<-"+r.RemoteAddr+"> results:", num)

	//load external files
	css, err := ioutil.ReadFile("ui/search.css")
	if err != nil {
		panic(err)
	}

	//add the page
	w.Write([]byte("<html><head><title>Distru :: Searching " + searchTerm + "</title><style type=\"text/css\">"))
	w.Write(css)
	w.Write([]byte("</style></head><body><div class=\"searchterm\">" + strconv.Itoa(num) + " results for <strong>" + searchTerm + "</strong></div>"))

	for i := range results {
		w.Write([]byte("<div class=\"results\"><a href=\"" + results[i].Link + "\">test</a></div>"))
	}
	w.Write([]byte("</body></html>"))
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
	w.Write([]byte("</script></head><body><div class = \"version\">" + Version + "<div class = \"name\">Distru</div><input type=\"text\" onkeydown=\"searchThis();\" onkeypress=\"isEnter(event);\" id=\"search\" class=\"search\" placeholder=\"Search freely\"/></body></html>"))
}
