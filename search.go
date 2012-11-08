package main

import (
	"time"
)

//Conf.Search returns the total number of results, and a []*page containing at most maxResults number of results.
func (conf *config) Search(terms []string) []*page {
	index := conf.Idx

	//Request indexes from all resources,
	//and trust their results.
	for i := range conf.Resources {
		index.MergeRemote(conf.Resources[i], true, conf.ResTimeout)
	}

	bareresults := make(map[string]*page, 0)
	for k, v := range conf.Idx.Cache {
		//For each page in the cache, we'll
		//check for the exact presence of each
		//search term.
		for i := range terms {
			_, isPresent := v.WordCount[terms[i]]
			if isPresent {
				//Refresh the timestamp on this page,
				//because it matched the search.
				v.Time = time.Now()
				bareresults[k] = v
			}
		}
	}
	for k, v := range index.Sites {
		for kk, vv := range v.Pages {
			//For each term, we get the number and presence
			//of the word for a particular page. The number
			//is currently discarded, because we can't rank
			//the relevance of pages.
			for i := range terms {
				_, isPresent := vv.WordCount[terms[i]]
				if isPresent {
					//Stamp the result with the
					//current time, for when it
					//is included in the chache.
					vv.Time = time.Now()
					bareresults[k+kk] = vv
				}
			}
		}
	}
	//The results should be sorted and refined before being returned, to improve the final search results.

	results := make([]*page, len(bareresults))
	for k, v := range bareresults {
		//Cache the results, but do not overwrite old results.
		_, isPresent := conf.Idx.Cache[k]
		if !isPresent {
			conf.Idx.Cache[k] = v
		}
		results = append(results, v)
	}
	return results
}
