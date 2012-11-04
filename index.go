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
	"time"
)

const (
	BotName = "Distru"
)

type Index struct {
	Sites map[string]*site //A map of fully indexed webpages.
	Queue chan string      `json:"-"` //The channel which controls Indexers
}

type site struct {
	Time  time.Time        //The time when the site finished indexing
	Pages map[string]*page //Nonordered map of pages on the server
	Links []string         //List of all unique links collected from all pages on the site
}

type page struct {
	Link      string         //The fully qualified link to this page
	WordCount map[string]int `json:"-"` //Temporary storage for the content of the page
}

//Index.Search
func (index *Index) Search(terms []string, maxResults int) []*page {
	results := make([]*page, 0)
	for _, v := range index.Sites {
		for _, vv := range v.Pages {
			//For each term, we get the number and presence
			//of the word for a particular page. The number
			//is currently discarded, because we can't rank
			//the relevance of pages.
			for i := range terms {
				_, isPresent := vv.WordCount[terms[i]]
				if isPresent {
					results = append(results, vv)
				}
			}
		}
	}
	return results
}

//Index.MergeRemote makes a raw distru request for the JSON encoded index of the given site, (which must have a full URI.) It will not overwrite local sites with remote ones unless trustNew is true. It returns nil if successful, or returns an error if the remote site could not be reached, or produced an invalid index.
func (index *Index) MergeRemote(remote string, trustNew bool) error {
	//Dial the connection here.
	conn, err := net.Dial("tcp", remote+":9049")
	if err != nil {
		return err
	}
	defer conn.Close()

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

	//Create a new decoder for reading the JSON
	//directly off the wire.
	dec := json.NewDecoder(r)

	var remoteIndex Index

	//Decode into the remoteIndex object.
	err = dec.Decode(&remoteIndex)
	if err != nil {
		return err
	}

	isPresent := false
	for k, v := range remoteIndex.Sites {
		//If the site in the local index is not present, or if
		//the remote index is trusted, *and* newer than the
		//local one, add the remote site to the local index.
		_, isPresent = index.Sites[k]
		if !isPresent || (trustNew && index.Sites[k].Time.Before(v.Time)) {
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
		newSite := newSite(target)
		if newSite == nil {
			//If we got an error for some reason,
			log.Println("indexer> failed to add \"" + target + "\"")
			//discard it and continue.
			continue
		}
		//Update the target site.
		index.Sites[target] = newSite
		log.Println("indexer> added \"" + target + "\"")
	}
}

func newSite(target string) *site {
	target = "http://" + target

	//Create an http.Client to control the webpage requests.
	client := http.Client{}
	//Use robotstxt to get the search engine permission.
	rperm, _ := getRobotsPermission(target)

	//Check if we are allowed to access /
	if !rperm.Test("/") {
		//If we aren't, return empty.
		return nil
	}

	pages := make(map[string]*page)
	newPage, tree, links := getPage(target, "/", client)
	if newPage == nil {
		//If we didn't get the root page properly,
		//return nil.
		return nil
	}
	//If we did get it, then continue normally.
	pages["/"] = newPage
	//Grab the root page first, then we're going to build on the tree.
	//We'll loop until there are no more unresolved pages. Then we'll
	//set isFinished to true, and break the loop.
	isFinished := false
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
			newPage, newTree, newLinks := getPage(target, k, client)
			if newPage == nil {
				//If we got a nil response from getPage,
				//then continue and drop this page
				continue
			}
			//If we got a good response, then put it in the map.
			pages[k] = newPage

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

	site := &site{
		Time:  time.Now(),
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
