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

	//add the <html> element
	fmt.Fprint(w, "<html>")
	//add the <head> element
	fmt.Fprint(w, "<head>")
	//add the title of the document and the search term
	fmt.Fprintf(w, "<title>Distru :: Searching \"%s\"</title>", searchTerm)

	//add the shortcuticon
	//BUG: THIS DOESNT WORK YET
	fmt.Fprint(w, "<link rel=\"shortcut icon\" href=\"img/icon_16.png\">")

	//add the stylesheet
	fmt.Fprint(w, "<style type=\"text/css\">")
	file, err := ioutil.ReadFile("ui/search.css")
	if err != nil {
		panic(err)
	}
	fmt.Fprint(w, string(file))
	fmt.Fprint(w, "</style>")
	//close the <head> element
	fmt.Fprint(w, "</head>")
	//add the <body> element
	fmt.Fprint(w, "<body>")
	//display the search term at the top
	fmt.Fprintf(w, "<div class=\"searchterm\">%d results for <strong>%s</strong></div>", numResults, searchTerm)

	//TODO: SEARCH HERE.
	//this is a temporary example of what searches will look like
	fmt.Fprint(w, "<div class=\"results\">test</div>")
	fmt.Fprint(w, "<div class=\"results\">test2</div>")

	//close the <body> element
	fmt.Fprint(w, "</body>")
	//close the <html> element
	fmt.Fprint(w, "</html>")
}

func frontpageHandler(w http.ResponseWriter, r *http.Request) {
	//add the <html> element
	fmt.Fprint(w, "<html>")
	//add the <head> element
	fmt.Fprint(w, "<head>")
	//add the title of the document
	fmt.Fprint(w, "<title>Distru :: Search Freely</title>")

	//add the shortcuticon
	//BUG: THIS DOESNT WORK YET
	fmt.Fprint(w, "<link rel=\"shortcut icon\" href=\"img/icon_16.png\">")

	//add the stylesheet
	fmt.Fprint(w, "<style type=\"text/css\">")
	file, err := ioutil.ReadFile("ui/index.css")
	if err != nil {
		panic(err)
	}
	fmt.Fprint(w, string(file))
	fmt.Fprint(w, "</style>")
	//add the javascript file
	fmt.Fprint(w, "<script type=\"text/javascript\">")
	file, err = ioutil.ReadFile("ui/distru.js")
	if err != nil {
		panic(err)
	}
	fmt.Fprint(w, string(file))
	fmt.Fprint(w, "</script>")
	//close the <head> element
	fmt.Fprint(w, "</head>")
	//add the <body> element
	fmt.Fprint(w, "<body>")
	//add the name that hovers above the search bar
	fmt.Fprint(w, "<div class = \"name\">Distru</div>")
	//add the form
	fmt.Fprint(w, "<input type=\"text\" onkeydown=\"searchThis();\" onkeypress=\"isEnter(event);\" id=\"search\" class=\"search\" placeholder=\"Search freely\"/>")
	//close the <body> element
	fmt.Fprint(w, "</body>")
	//close the <html> element
	fmt.Fprint(w, "</html>")
}
