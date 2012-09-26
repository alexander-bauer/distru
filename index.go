package main

//import "io/ioutil"

type Index struct {
	Sites []site //list of indexed webpages
}

type site struct {
	URL   string     //the link that identifies this Block
	Pages []sitePage //nonordered list of pages and their data on the server
}

type sitePage struct {
	Path    string   //path to page on the webserver (relative to root page)
	Links   []string //list of hyperlinks on the page
	Content string   //the content, temporarily replacing word lists
}

//NewIndex is a constructor for the Index struct
func NewIndex() *Index {
	
	//get peer list here TODO
	
	peerList := []string{"localhost","example.com"}
	peerSites := make([]site, len(peerList))
	
	for i := range peerList {
		peerSites[i] = *newSite(peerList[i])
	}
	
	index := Index{
		Sites: peerSites,
	}
	return &index
}

/*
func (index *Index) save(path string) error {
	binary := byte(index)
	return ioutil.WriteFile(path, index, 0600)
}*/


func newSite(url string) *site {
	pages := []sitePage{*newSitePage(url)} //make an array of length 1
	//by scraping the site's page
	//TODO this should build the whole tree

	site := site{
		URL:   url,
		Pages: pages,
	}
	return &site
}

//sitePage constructor, which scrapes a single webpage
func newSitePage(url string) *sitePage {
	body := fetch(url)      //get the body of the webpage
	allLinks := getLinks(body) //collect links, but
	links := getExternalLinks(allLinks) //get only the external links

	//the wordlist should be added here, but that function doesn't exist yet
	//TODO

	page := sitePage{
		Path:    url,
		Links:   links,
		Content: body,
	}
	return &page
}
