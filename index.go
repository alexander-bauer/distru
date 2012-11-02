package main

import (
	"bufio"
	"encoding/gob"
	"encoding/json"
	"github.com/temoto/robotstxt.go"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"strings"
	"bytes"
)

const (
	BotName = "Distru"
)

type Index struct {
	Sites map[string]site //A map of fully indexed webpages.
	Queue chan string     `json:"-"` //The channel which controls Indexers
}

type site struct {
	Pages map[string]*page //Nonordered map of pages on the server
	Links []string         //List of all unique links collected from all pages on the site
}

type page struct {
	Content string //Temporary storage for the content of the page
}

//Index.MergeRemote makes a raw distru request for the JSON encoded index of the given site, (which must have a full URI.) It will not overwrite local sites with remote ones. It returns nil if successful, or returns an error if the remote site could not be reached, or produced an invalid index.
func (index *Index) MergeRemote(remote string) error {
	//Dial the connection here.
	conn, err := net.Dial("tcp", remote)
	if err != nil {
		return err
	}
	//Initialize a new reader and writer.
	r, w := bufio.NewReader(conn), bufio.NewWriter(conn)
	_, err = w.WriteString(GETJSON) //Request the JSON-encoded index.
	if err != nil {
		return err
	}

	err = w.Flush() //Flush the writer to the connection.
	if err != nil {
		return err
	}

	resp, err := r.ReadBytes('\n') //Read the response.
	if err != nil {
		return err
	}

	remoteIndex := &Index{}
	err = json.Unmarshal(resp, remoteIndex) //Marshal into an index object
	if err != nil {
		return err
	}
	isPresent := false

	for k, v := range remoteIndex.Sites {
		//If the local index contains the site already,
		//don't overwrite it.
		_, isPresent = index.Sites[k]
		if !isPresent {
			//Otherwise, add the remote index's site
			//to the local index.
			index.Sites[k] = v
		}
		//Repeat until we've gone through all of the
		//values in remoteIndex.
	}
	return nil
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

//MaintainIndex launches a number of goroutines which handle indexing of sites in sequence. It sets index.Queue to a channel into which target urls should be placed. When a new string is added to the returned chan, one of the next non-busy indexer will remove it from the chan and index it, and add the contents to the passed index. It will then forget about that site.
//To remove a site from the index, use delete(index.Sites, urlstring). To shut down the indexers, close() index.Queue.
func MaintainIndex(index *Index, numIndexers int) {
	//First, we're going to make the channel of pending sites.
	index.Queue = make(chan string)

	//Next, we're going to launch numIndexers amount of Indexers.
	for i := 0; i < numIndexers; i++ {
		go Indexer(index, index.Queue)
	}
}

func Indexer(index *Index, pending <-chan string) {
	for target := range pending {
		//Update the target site.
		index.Sites[target] = newSite(target)
		log.Println("indexer> added \"" + target + "\"")
	}
}

func newSite(target string) site {
	target = "http://" + target
	//Initialize an empty tree and set isFinished to false.
	tree := make(map[string]struct{})
	links := make(map[string]struct{})
	isFinished := false

	//Create an http.Client to control the webpage requests.
	client := http.Client{}
	//Use robotstxt to get the search engine permission.
	rperm, _ := getRobotsPermission(target)

	//Check if we are allowed to access /
	if !rperm.Test("/") {
		//If we aren't, return empty.
		return site{}
	}

	pages := make(map[string]*page)
	pages["/"], tree, links = getPage(target, "/", client)
	//Grab the root page first, then we're going to build on the tree.
	//We'll loop until there are no more unresolved pages. Then we'll
	//set isFinished to true, and break the loop.
	for isFinished == false {
		//We set isFinished to true here. If we're not actually
		//finished, the following loop will set it to false.
		isFinished = true
		for k, _ := range tree {
			if pages[k] != nil || !rperm.Test(k) {
				//If the page has been indexed already,
				//or if we're not allowed to access it,
				//ignore it.
				continue
			}
			//Otherwise, set isFinished to false, because we will
			//need at least one more iteration.
			isFinished = false
			//Then we index the page and grab the new tree.
			newTree := make(map[string]struct{})
			newLinks := make(map[string]struct{})
			pages[k], newTree, newLinks = getPage(target, k, client)

			//Then we put all of the new values into the old maps,
			for kk, vv := range newTree {
				tree[kk] = vv
			}
			for kk, vv := range newLinks {
				links[kk] = vv
			}
			//and start the loop over again.
		}
	}
	linkArray := make([]string, 0, len(links))
	for k, _ := range links {
		linkArray = append(linkArray, k)
	}

	site := site{
		Pages: pages,
		Links: linkArray,
	}
	return site
}

func getRobotsPermission(target string) (*robotstxt.Group, error) {
	//We're going to define a routine with which to fail.
	fail := func(err error) (*robotstxt.Group, error) {
		//Since we're failing here when there is no file available,
		//craft a stand-in one to be parsed instead.
		robots, _ := robotstxt.FromBytes([]byte("User-agent: *\nAllow: /"))
		return robots.FindGroup(BotName), err
	}
	//Use robotstxt here.
	resp, err := http.Get(target + "/robots.txt")
	if err != nil {
		return fail(err)
	}
	defer resp.Body.Close()
	robots, err := robotstxt.FromResponse(resp)
	if err != nil {
		return fail(err)
	}
	group := robots.FindGroup(BotName)
	if group == nil {
		//BUG(DuoNoxSol): This does not raise a real error.
		return fail(nil)
	}
	return group, nil
}

//getPage is a complex constructor for the page object. It appends path to target in order to get the target webpage. It then uses http.Get to get the body of that webpage, which it then uses regexp to scrape for links. Those links are sorted into internal and external. The internal links are resolved to be absolute (internal) links on the webserver, and then returned, without duplicates, as a map[string]struct{}. All unique external links on the page are returned in the second map[string]struct{}.
func getPage(target, path string, client http.Client) (*page, map[string]struct{}, map[string]struct{}) {
	//Parse the target URI, return empty if it fails.
	accessURI, err := url.ParseRequestURI(target + path)
	if err != nil {
		return &page{}, nil, nil
	}

	//Get the content of the webpage via HTTP, using the
	//existing http.Client, and return blank if it fails.
	resp, err := client.Get(accessURI.String())
	if err != nil {
		return &page{}, nil, nil
	}
	defer resp.Body.Close()
	//Get the body of the request as a []byte.
	b, err := ioutil.ReadAll(resp.Body)
	//Convert to string real quick.
	body := string(b)

	//Now we're going to move on to parsing the links.
	pattern, err := regexp.Compile("href=['\"]?([^'\" >]+)")
	if err != nil {
		return &page{}, nil, nil
	}

	//Use pattern matching to find all link tags on the page,
	//and put them in array.
	tags := pattern.FindAllStringSubmatch(body, -1)

	//Now parse them into a list of actual links.
	//We're going to separate the internal and external
	//links in the same step.
	internalLinks := make(map[string]struct{}, len(tags)) //Reserve space
	externalLinks := make(map[string]struct{}, len(tags)) //for len(tags) items.

	for i := range tags {
		//tags is an array containing both the "href=" and the link
		link := tags[i][1] //so we take only the link element

		if !strings.Contains(link, "http://") && !strings.Contains(link, "https://") {
			//If the string doesn't contain http://,
			//resolve it to an absolute 
			internalLinks[join(path, link)] = struct{}{}
		} else {
			//If the string directs to this site (with http://)
			//then put it in internal links
			if strings.HasPrefix(link, target) {
				//(but trim the website name
				internalLinks[join(path, link[len(target):])] = struct{}{}
				//and jump back to the beginning of the for,)
				continue
			}
			//otherwise, put it in externals.
			externalLinks[link] = struct{}{}
		}
	}

	//the wordlist should be added here, but that function doesn't exist yet
	//TODO
	
	//only lowercase letters!
	b = bytes.ToLower(b)
	
	//Compile the pattern for stripping HTML
	p, err := regexp.Compile("<([^>]*)>|\n|\u0009")
	if err != nil {
		return &page{}, nil, nil
	}
	//apply the pattern
	body = string(p.ReplaceAll(b, []byte("")))
	
	
	return &page{
		Content: body,
	}, internalLinks, externalLinks
}

func join(source, target string) string {
	if path.IsAbs(target) {
		return target
	}
	return path.Join(path.Dir(source), target)
}
