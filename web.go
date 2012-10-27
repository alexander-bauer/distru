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
	log.Println("Starting webserver.")
	http.ListenAndServe(":9048", nil)

}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	searchTerm := r.URL.Path[len("/search/"):]
	log.Println("<-" + r.RemoteAddr + "> searching \"" + searchTerm + "\"")
	fmt.Fprint(w, searchTerm)
}

func frontpageHandler(w http.ResponseWriter, r *http.Request) {

	//add the <html> element
	fmt.Fprint(w, "<html>")

	//add the <head> element
	fmt.Fprint(w, "<head>")

	//add the stylesheet
	fmt.Fprint(w, "<style type=\"text/css\">")
	file, err := ioutil.ReadFile("style.css")
	if err != nil {
		panic(err)
	}
	fmt.Fprint(w, string(file))
	fmt.Fprint(w, "</style>")

	//add the javascript file
	fmt.Fprint(w, "<script type=\"text/javascript\">")
	fmt.Fprint(w, "function searchThis() {if (event.keyCode == 13) window.location = '/search/'+document.getElementById('search').value;}")
	fmt.Fprint(w, "</script>")

	//close the <head> element
	fmt.Fprint(w, "</head>")

	//add the <body> element
	fmt.Fprint(w, "<body>")

	//add the name that hovers above the search bar
	fmt.Fprint(w, "<div class = \"name\">Distru</div>")

	//add the form
	fmt.Fprint(w, "<input type=\"text\" onkeydown=\"searchThis()\" id=\"search\" class=\"search\" placeholder=\"Search freely\"/>")

	//close the <body> element
	fmt.Fprint(w, "</body>")

	//close the <html> element
	fmt.Fprint(w, "</html>")

}
