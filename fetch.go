package main

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
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
	if err != nil {                           //if there's an error,
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

func getInternalLinks(links []string, s string) []string {
	in := [999]string{} //eventually we need to change the "999"
	ixcount := 0

	for i := range links {
		if strings.Contains(s, links[i]) && !strings.Contains(links[i], "http://") {
			in[ixcount] = links[i]
			ixcount++
		}
	}

	//this part gets rid of the extra 900 or so spaces that are there
	internal := make([]string, ixcount)
	for i := 0; i < ixcount; i++ {
		internal[i] = in[i]
	}
	return (internal)
}

func getExternalLinks(links []string, s string) []string {
	ex := [999]string{} //eventually we need to change the "999"
	excount := 0

	for i := range links {
		if strings.Contains(links[i], "http://") {
			ex[excount] = links[i]
			excount++
		}
	}

	//this part gets rid of the extra 900 or so spaces that are there
	external := make([]string, excount)
	for i := 0; i < excount; i++ {
		external[i] = ex[i]
	}
	return (external)
}

func main() {
	s := fetch(os.Args[1])
	links := getLinks(s)

	internal := getInternalLinks(links, s)
	external := getExternalLinks(links, s)

	print("Internal Links:\n")
	for i := range internal {
		print(internal[i], "\n")
	}

	print("\n")

	print("External Links:\n")
	for i := range external {
		print(external[i], "\n")
	}

	os.Exit(0)
}
