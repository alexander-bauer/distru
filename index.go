package main

import (
	"encoding/gob"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"strings"
)

type Index struct {
	Sites map[string]site //A map of fully indexed webpages.
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

//MaintainIndex launches a number of goroutines which handle indexing of sites in sequence. It returns a chan into which target urls should be placed. When a new string is added to the returned chan, one of the next non-busy indexer will remove it from the chan and index it, and add the contents to the passed index. It will then forget about that site.
//To remove a site from the index, use delete(index.Sites, urlstring).
func MaintainIndex(index *Index, numIndexers int) chan<- string {
	//First, we're going to make the channel of pending sites.
	pending := make(chan string)

	//Next, we're going to launch numIndexers amount of Indexers.
	for i := 0; i < numIndexers; i++ {
		go Indexer(index, pending)
	}

	return pending
}

func Indexer(index *Index, pending <-chan string) {
	for target := range pending {
		//Update the target site.
		index.Sites[target] = newSite(target)
		log.Println("indexer> added \"" + target + "\"")
	}
}

type site struct {
	URL   string           //domain or IP that identifies this Block
	Pages map[string]*page //nonordered map of pages on the server
}

func newSite(target string) site {
	//Initialize an empty tree and set isFinished to false.
	tree := []string{}
	isFinished := false

	pages := make(map[string]*page)
	pages["/"], tree = getPage(target, "/")
	//Grab the root page first, then we're going to build on the tree.
	//We'll loop until there are no more unresolved pages. Then we'll
	//set isFinished to true, and break the loop.
	for isFinished == false {
		//We set isFinished to true here. If we're not actually
		//finished, the following loop will set it to false.
		isFinished = true
		for i := range tree {
			if pages[tree[i]] != nil {
				//If the page has been indexed already,
				//ignore it.
				continue
			}
			//Otherwise, set isFinished to false, because we will
			//need at least one more iteration.
			isFinished = false
			//Then we index the page and grab the new tree.
			newTree := []string{}
			pages[tree[i]], newTree = getPage(target, tree[i])

			//Then we append the new tree to the old one,
			tree = append(tree, newTree...)
			//and start the loop over again.
		}
	}

	site := site{
		URL:   target,
		Pages: pages,
	}
	return site
}

type page struct {
	Path  string   //path to page on the webserver (relative to root page)
	Links []string //list of external links on the page
}

//getPage is a complex constructor for the page object. It appends path to target in order to get the target webpage. It then uses http.Get to get the body of that webpage, which it then uses regexp to scrape for links. Those links are sorted into internal and external. The external links are put into the Links element of the page structure. The internal links are resolved to be absolute (internal) links on the webserver, and then returned, possibly with duplicates, as a []string, in which every element is true.
func getPage(target string, path string) (*page, []string) {
	//Parse the target URI, return empty if it fails.
	accessURI, err := url.ParseRequestURI(target + path)
	if err != nil {
		//Prepend http:// permanently
		target = "http://" + target
		accessURI, err = url.ParseRequestURI(target + path)
		if err != nil {
			return &page{}, nil
		}
	}

	//Get the content of the webpage via HTTP, return blank if it fails.
	resp, err := http.Get(accessURI.String())
	if err != nil {
		return &page{}, nil
	}
	defer resp.Body.Close()
	//Get the body of the request as a []byte.
	b, err := ioutil.ReadAll(resp.Body)
	//Convert to string real quick.
	body := string(b)

	//Now we're going to move on to parsing the links.
	pattern, err := regexp.Compile("href=['\"]?([^'\" >]+)")
	if err != nil {
		return &page{}, nil
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

		if !strings.Contains(link, "http://") && !strings.Contains(link, "https://") {
			//If the string doesn't contain http://,
			//resolve it to an absolute 
			internalLinks = append(internalLinks, join(path, link))
		} else {
			//If the string directs to this site (with http://)
			//then put it in internal links
			if strings.HasPrefix(link, target) {
				//(but trim the website name
				internalLinks = append(internalLinks, join(path, link[len(target):]))
				//and jump back to the beginning of the for,)
				continue
			}
			//otherwise, put it in externals.
			externalLinks = append(externalLinks, link)
		}
	}

	//the wordlist should be added here, but that function doesn't exist yet
	//TODO

	return &page{
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
