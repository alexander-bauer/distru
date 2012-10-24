package main

import (
	"fmt"
	"log"
	"net/http"
	"io/ioutil"
)

func ServeWeb() {
	http.HandleFunc("/", searchHandler)
	log.Println("Starting webserver.")
	http.ListenAndServe(":9048", nil)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	//add the <html> element
	s := "<html>"
	fmt.Fprint(w, s)
	
	//add the <head> element
	s = "<head>"
	fmt.Fprint(w, s)
	
	//add the stylesheet
	s = "<style type='text/css'>"
	fmt.Fprint(w, s)
	file, err := ioutil.ReadFile("webui/style.css")
		if err != nil { panic(err) }
	s = string(file)
	fmt.Fprint(w, s)
	s = "</style>"
	fmt.Fprint(w, s)
	
	//add the jquery
	s = "<script type='text/javascript'>"
	fmt.Fprint(w, s)
	file, err = ioutil.ReadFile("webui/jquery.js")
		if err != nil { panic(err) }
	s = string(file)
	fmt.Fprint(w, s)
	s = "</script>"
	fmt.Fprint(w, s)
	
	//add the javascript file
	s = "<script type='text/javascript'>"
	fmt.Fprint(w, s)
	file, err = ioutil.ReadFile("webui/common.js")
		if err != nil { panic(err) }
	s = string(file)
	fmt.Fprint(w, s)
	s = "</script>"
	fmt.Fprint(w, s)
		
	//close the <head> element
	s = "</head>"
	fmt.Fprint(w, s)
	
	//add the body of index.html
	file, err = ioutil.ReadFile("webui/index.html")
		if err != nil { panic(err) }
	s = string(file)
	fmt.Fprint(w, s)

	//close the <html> element
	s = "</html>"
	fmt.Fprint(w, s)
	
}
