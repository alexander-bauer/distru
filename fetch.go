package main

import (
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
)

//fetch a webpage from url (without http://)
func fetch(url string) string {
	resp, err := http.Get("http://" + url) //make the request
	if err != nil {            //if there's an error,
		os.Exit(1) //then exit with error 1
	}
	defer resp.Body.Close()                //(not sure what this does)
	body, err := ioutil.ReadAll(resp.Body) //get the body of the request in []byte form
	content := string(body)                //convert to string
	return (content)
}

func getLinks(html string) [][]string {
	tags, tagErr := regexp.Compile("href=['\"]?([^'\" >]+)")
	if tagErr != nil {
		os.Exit(1)
	}
	
	/*href, hrefErr := regexp.Compile("\\s*(?i)href\\s*=\\s*(\"([^\"]*\")|'[^']*'|([^'\">\\s]+))")
	if hrefErr != nil {
		os.Exit(1)
	}*/
	
	return (tags.FindAllStringSubmatch(html, -1))
}

func scrapeurl(content string) []string {
	return(nil)
	}

func main() {
	s := fetch(os.Args[1])
	links := getLinks(s)
	
	for i := range links {
		print(links[i][1], "\n")
	}
}
