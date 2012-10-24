package main

import (
	"encoding/gob"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

type Index struct {
	Sites []site //list of indexed webpages
}

//Index.Gob uses encoding/gob to write a binary representation of itself to the specified io.writer. This can be used to pass indexes across Conn objects.
func (index *Index) Gob(w io.Writer) {
	gob.NewEncoder(w).Encode(Idx)
}

//Index.JSON creates a JSON-encoded (encoding/json) and tab indented string from the parent index.
func (index *Index) JSON() string {
	//We're going to marshal the parent index here with tab indentation
	b, err := json.MarshalIndent(index, "", "\t")
	if err != nil {
		return "" //return a blank string if there's an error
	}

	//Then return the []byte as converted to a string.
	return string(b)
}

//NewIndex is a constructor for the Index struct
func NewIndex() *Index {

	//get peer list here TODO

	peerList := []string{"localhost", "uppit.us", "example.com"}
	peerSites := make([]site, len(peerList))

	for i := range peerList {
		peerSites[i] = *newSite(peerList[i])
	}

	index := Index{
		Sites: peerSites,
	}
	return &index
}

type site struct {
	URL   string //domain or IP that identifies this Block
	Pages []page //nonordered list of pages and their data on the server
	Tree  []string
}

func newSite(target string) *site {
	pages := []page{*newPage(target, "/")} //make an array of length 1
	//by scraping the site's page
	//TODO this should build the whole tree

	site := site{
		URL:   target,
		Pages: pages,
		Tree:  pages[0].Internals, //just grabbing /'s internal links
	}
	return &site
}

type page struct {
	Path      string   //path to page on the webserver (relative to root page)
	Internals []string //list of internal links on the page
	Externals []string //list of external links on the page
}

//newPage is the page constructor. It takes a target URL, (that being the base
//of the website, without trailing /,) and path to the target webpage. It
//fetches the webpage by combining the target and path, then scrapes the links
//from the body of the html. It determines whether each link is an internal or
//external link, and puts them in different arrays, then returns a page
//containing the resulting information. It returns an empty page if it
//encounters an error.
func newPage(target string, path string) *page {
	//Parse the target URI, return empty if it fails.
	accessURI, err := url.ParseRequestURI(target + path)
	if err != nil {
		accessURI, err = url.ParseRequestURI("http://" + target + path)
		if err != nil {
			return &page{}
		}
	}

	//Get the content of the webpage via HTTP, return blank if it fails.
	resp, err := http.Get(accessURI.String())
	if err != nil {
		return &page{}
	}
	defer resp.Body.Close()
	//Get the body of the request as a []byte.
	b, err := ioutil.ReadAll(resp.Body)
	//Convert to string real quick.
	body := string(b)

	//Now we're going to move on to parsing the links.
	pattern, err := regexp.Compile("href=['\"]?([^'\" >]+)")
	if err != nil {
		return &page{}
	}

	//Use pattern matching to find all link tags on the page,
	//and put them in array.
	tags := pattern.FindAllStringSubmatch(body, -1)

	//Now parse them into a list of actual links.
	//We're going to separate the internal and external
	//links in the same step.
	internalLinks := make([]string, 0, len(tags)) //length 0, reserve space
	externalLinks := make([]string, 0, len(tags)) //for len(tags) items

	for i := range tags {
		//tags is an array containing both the "href=" and the link
		link := tags[i][1] //so we take only the link element

		if !strings.Contains(link, "http://") {
			//If the string doesn't contain http://,
			//put it in the internal section
			internalLinks = append(internalLinks, link)
		} else {
			//otherwise, put it in externals
			internalLinks = append(externalLinks, link)
		}
	}

	//the wordlist should be added here, but that function doesn't exist yet
	//TODO

	return &page{
		Path:      path,
		Internals: internalLinks,
		Externals: externalLinks,
	}
}
