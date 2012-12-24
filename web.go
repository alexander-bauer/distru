package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var webDir string

func ServeWeb(conf *config) {
	http.HandleFunc("/search/", searchHandler)
	http.HandleFunc("/", frontpageHandler)
	webDir = conf.WebDir
	log.Println("Started webserver on port 9048.")
	http.ListenAndServe(":9048", nil)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	//get the search term and save it as searchTerm
	searchTerms := r.URL.Path[len("/search/"):]
	log.Println("<-" + r.RemoteAddr + "> searching \"" + searchTerms + "\"")

	//Perform the search.
	results := Conf.Search(strings.Split(searchTerms, " "))

	log.Println("<-"+r.RemoteAddr+"> results:", len(results))

	//load external files
	css, err := ioutil.ReadFile(webDir + "/search.css")
	if err != nil {
		panic(err)
	}
	parseJS, err := ioutil.ReadFile(webDir + "/parse.js")
	if err != nil {
		panic(err)
	}

	searchJS, err := ioutil.ReadFile(webDir + "/search.js")
	if err != nil {
		panic(err)
	}

	//add the page
	w.Write([]byte("<html><head><title>Distru :: Searching " + searchTerms + "</title><div class = \"version\">" + Version + "</div><style type=\"text/css\">"))
	w.Write(css)
	w.Write([]byte("</style></head><body><div class =\"holder\"><div class=\"searchterm\">" + strconv.Itoa(len(results)) + " results for <span id=\"term\"><strong>" + searchTerms + "</strong></span><a href=\"/\" class=\"back\">Distru</a><input type=\"text\" onkeydown=\"searchThis();\" onkeypress=\"isEnter(event);\" id=\"search\" class=\"search\" placeholder=\"Search freely\"/></div></div><div id=\"blank\"></div><script type=\"text/javascript\">"))
	w.Write(parseJS)
	w.Write(searchJS)
	w.Write([]byte("</script>"))

	for i := range results {
		//get url and remove the http://
		url := results[i].Link[len("http://"):]
		//if the url has a "/" at the end, remove it
		if strings.HasSuffix(url, "/") {
			url = url[:len(url)-1]
		}

		w.Write([]byte("<div class =\"holder\"><a href=\"" + results[i].Link + "\"><div class=\"results\">" + results[i].Title + "<div class =\"description\">" + results[i].Description + "</div><div class=\"url\">" + url + "</div></div></a></div>"))

	}
	w.Write([]byte("</body></html>"))
}

func frontpageHandler(w http.ResponseWriter, r *http.Request) {
	//load external files
	css, err := ioutil.ReadFile(webDir + "/index.css")
	if err != nil {
		panic(err)
	}
	javascript, err := ioutil.ReadFile(webDir + "/search.js")
	if err != nil {
		panic(err)
	}

	//add the page
	w.Write([]byte("<html><head><title>Distru :: Search Freely</title><style type=\"text/css\">"))
	w.Write(css)
	w.Write([]byte("</style><script type=\"text/javascript\">"))
	w.Write(javascript)
	w.Write([]byte("</script></head><body><div class = \"version\">" + Version + "</div><div class = \"name\">Distru</div><input type=\"text\" onkeydown=\"searchThis();\" onkeypress=\"isEnter(event);\" id=\"search\" class=\"search\" placeholder=\"Search freely\" autofocus/></body></html>"))
}
