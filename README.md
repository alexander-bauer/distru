# Distru

**Distru** is a still **in-progress** tool, written in [Go](http://golang.org), designed to act as a distributed search engine for [Hyperboria](https://projectmeshnet.org). Distru is not a single "search engine server," but rather a small program designed to maintain information about nearby servers, so that Hyperboria nodes can quickly share discovered sites.

## Concept

Search engines on the *Old Internet* crawl the internet, searching for websites, and index content. A search you make on an *Old Internet* search engine is simply a query to the engine's database. If the engine's servers go down, or your connection to them is cut off, then you're stuck. That sort of search engine *can't* be distributed.

Distru does not maintain a single or small number of databases. Instead, every node running distru "indexes" the content of its neighboring nodes. In addition, it temporarily indexes sites that it discovers. (For example, when one might search "reddit onmesh," and discover [uppit](http://uppit.us), that site is indexed for several hours or a day.)

When Distru makes a search request, it sends index queries to peers and indexed sites, which respond with their entire indexes. Distru may then send queries to all of the sites in the newly built index, depending on the size of the index.

Then, it consults the word lists in the index, and displays sites with words similar to the ones searched for. The top result sites are added to the temporary index, and the rest of the collected index is discarded.

### Indexing

"Indexing" a site has a very specific meaning. Distru has the ability to download webpages, and then analyze their content. This entails scraping all hyperlinks and site-relative links from the page, and building a list of non-common words in it. (*Coming soon.*) Distru can then build a "tree" of the entire site, based on the links between its own webpages.

When a site's tree of content is known, and all of the hyperlinks from those pages are known, and there is a word list for each of those pages, Distru has indexed that site. Distru's index is built of a list of indexed sites.

## To Use

Distru is written in Go, which is terribly easy to compile. As Distru nears a more-finished state, binaries will be available for download, of course. Until then, it can be downloaded and executed as follows.

```
sudo apt-get install golang

git clone git@github.com:SashaCrofter/distru.git
cd distru
go build
```

(*If you do not have a GitHub account, you may need to use* `https://github.com/SashaCrofter/distru.git` *instead.*) After using `go build`, the file `distru` will be executable and in the current directory. In the future, you may allow it to make an index, or run certain commands, but currently you can only scrape webpages for URLs. 
