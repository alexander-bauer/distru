package main

import (
	"fmt"
	"log"
	"net/http"
	"io/ioutil"
)

func ServeWeb() {
	http.HandleFunc("/search/", searchHandler)
	http.HandleFunc("/", frontpageHandler)
	log.Println("Starting webserver.")
	http.ListenAndServe(":9048", nil)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	
}

func frontpageHandler(w http.ResponseWriter, r *http.Request) {
	//add the <html> element
	fmt.Fprint(w, "<html>")
	
	//add the <head> element
	fmt.Fprint(w, "<head>")
	
	//add the stylesheet
	fmt.Fprint(w, "<style type='text/css'>")
	file, err := ioutil.ReadFile("webui/style.css")
		if err != nil { panic(err) }
	fmt.Fprint(w, string(file))
	fmt.Fprint(w, "</style>")
	
	//add the jquery
	fmt.Fprint(w, "<script type='text/javascript'>")
	file, err = ioutil.ReadFile("webui/jquery.js")
		if err != nil { panic(err) }
	fmt.Fprint(w, string(file))
	fmt.Fprint(w, "</script>")
	
	//add the javascript file
	fmt.Fprint(w, "<script type='text/javascript'>")
	file, err = ioutil.ReadFile("webui/common.js")
		if err != nil { panic(err) }
	fmt.Fprint(w, string(file))
	fmt.Fprint(w, "</script>")
		
	//close the <head> element
	fmt.Fprint(w, "</head>")
	
	//add the body of index.html
	file, err = ioutil.ReadFile("webui/index.html")
		if err != nil { panic(err) }
	fmt.Fprint(w, string(file))

	//close the <html> element
	fmt.Fprint(w, "</html>")
	
}
