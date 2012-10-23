package main

import (
	"encoding/json"
)

type Index struct {
	Sites []site //list of indexed webpages
}

//Index.JSON creates a JSON-encoded and tab indented string from the parent index.
func (index *Index) JSON() string {
	//We're going to marshal the parent index here with tab indentation
	b, err := json.MarshalIndent(index, "", "\t")
	if err != nil {
		return "" //return a blank string if there's an error
	}
	
	//Then return the []byte as converted to a string.
	return string(b)
}

type site struct {
	URL   string //domain or IP that identifies this Block
	Pages []page //nonordered list of pages and their data on the server
	Tree  []string
}

type page struct {
	Path      string   //path to page on the webserver (relative to root page)
	Links     []string //list of hyperlinks on the page
	Internals []string //list of internal links on the page
	Content   string   //the content, temporarily replacing word lists
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

func newSite(url string) *site {
	pages := []page{*newPage(url, "/")} //make an array of length 1
	//by scraping the site's page
	//TODO this should build the whole tree

	site := site{
		URL:   url,
		Pages: pages,
		Tree:  pages[0].Internals, //just grabbing /'s internal links
	}
	return &site
}

//newPage is the site constructor, which scrapes a single webpage
//it takes the URL of a site (without trailing /,) and a directory path, such
//as / or /help.txt
//It returns the sitePage object, as well as an array of internal links.
func newPage(url string, path string) *page {
	body := fetch(url)                                //get the body of the webpage
	allLinks := getLinks(body)                        //collect links, but
	externalLinks := getExternalLinks(allLinks)       //get only the external links
	internalLinks := getInternalLinks(allLinks, body) //get only internal links

	//the wordlist should be added here, but that function doesn't exist yet
	//TODO

	return &page{
		Path:      path,
		Links:     externalLinks,
		Internals: internalLinks,
		Content:   body,
	}
}
