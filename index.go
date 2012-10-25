package main

import (
	"encoding/gob"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
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
		peerSites[i] = newSite(peerList[i])
	}

	index := Index{
		Sites: peerSites,
	}
	return &index
}

type site struct {
	URL   string //domain or IP that identifies this Block
	Pages []page //nonordered list of pages and their data on the server
	Tree  map[string]bool
}

func newSite(target string) site {

	root, tree := getPage(target, "/") //make an array of length 1
	//by scraping the site's page
	//TODO this should build the whole tree

	site := site{
		URL:   target,
		Pages: []page{root},
		Tree:  tree,
	}
	return site
}

type page struct {
	Path  string   //path to page on the webserver (relative to root page)
	Links []string //list of external links on the page
}

//getPage is a complex constructor for the page object. It appends path to target in order to get the target webpage. It then uses http.Get to get the body of that webpage, which it then uses regexp to scrape for links. Those links are sorted into internal and external. The external links are put into the Links element of the page structure. The internal links are resolved to be absolute (internal) links on the webserver, and then returned (without duplicates) as a map[string]bool, in which every element is true.
func getPage(target string, path string) (page, map[string]bool) {
	//Parse the target URI, return empty if it fails.
	accessURI, err := url.ParseRequestURI(target + path)
	if err != nil {
		//Prepend http:// permanently
		target = "http://" + target
		accessURI, err = url.ParseRequestURI(target + path)
		if err != nil {
			return page{}, nil
		}
	}

	//Get the content of the webpage via HTTP, return blank if it fails.
	resp, err := http.Get(accessURI.String())
	if err != nil {
		return page{}, nil
	}
	defer resp.Body.Close()
	//Get the body of the request as a []byte.
	b, err := ioutil.ReadAll(resp.Body)
	//Convert to string real quick.
	body := string(b)

	//Now we're going to move on to parsing the links.
	pattern, err := regexp.Compile("href=['\"]?([^'\" >]+)")
	if err != nil {
		return page{}, nil
	}

	//Use pattern matching to find all link tags on the page,
	//and put them in array.
	tags := pattern.FindAllStringSubmatch(body, -1)

	//Now parse them into a list of actual links.
	//We're going to separate the internal and external
	//links in the same step.
	internalLinks := make(map[string]bool, len(tags)) //length 0, reserve space
	externalLinks := make([]string, 0, len(tags))     //for len(tags) items

	for i := range tags {
		//tags is an array containing both the "href=" and the link
		link := tags[i][1] //so we take only the link element

		if !strings.Contains(link, "http://") && !strings.Contains(link, "https://") {
			//If the string doesn't contain http://,
			//resolve it to an absolute 
			internalLinks[join(path, link)] = true
		} else {
			//If the string directs to this site (with http://)
			//then put it in internal links
			if strings.HasPrefix(link, target) {
				//(but trim the website name
				internalLinks[join(path, link[len(target):])] = true
				//and jump back to the beginning of the for,)
				continue
			}
			//otherwise, put it in externals.
			externalLinks = append(externalLinks, link)
		}
	}

	//the wordlist should be added here, but that function doesn't exist yet
	//TODO

	return page{
		Path:  path,
		Links: externalLinks,
	}, internalLinks
}

func join(source, target string) string {
	if path.IsAbs(target) {
		return target
	}
	return path.Join(path.Dir(source), target)
}
