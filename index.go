package main

import (
	"bufio"
	"encoding/gob"
	"encoding/json"
	"github.com/temoto/robotstxt.go"
	"io"
	"log"
	"net"
	"net/http"
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
	WordCount map[string]int //Temporary storage for the content of the page
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
