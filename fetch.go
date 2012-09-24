package main

import (
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"net/url"
	"distru/index"
)

//fetch a webpage from url (without http://)
func fetch(path string) string {
	accessURI, err := url.ParseRequestURI(path)
	if err != nil {
		accessURI, err = url.ParseRequestURI("http://" + path)
		if err != nil {
			os.Exit(1)
		}
	}
	resp, err := http.Get(accessURI.String()) //make the request
	if err != nil {            //if there's an error,
		os.Exit(1) //then exit with error 1
	}
	defer resp.Body.Close()                //(not sure what this does)
	body, err := ioutil.ReadAll(resp.Body) //get the body of the request in []byte form
	content := string(body)                //convert to string
	return (content)
}

//return all URLs in href attributes of the given HTML
func getLinks(html string) []string {
	tags, tagErr := regexp.Compile("href=['\"]?([^'\" >]+)")
	if tagErr != nil {
		os.Exit(2)
	}
	
	links := tags.FindAllStringSubmatch(html, -1)
	
	linkTexts := make([]string, len(links))
	
	//We only want the second matched set [1], which does not contain 'http='
	for i := range links {
		linkTexts[i] = links[i][1]
	}
	return (linkTexts)
}

func main() {
	s := fetch(os.Args[1])
	links := getLinks(s)
	
	for i := range links {
		print(links[i], "\n")
	}
}
