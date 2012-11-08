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

	bareresults := make([]*page, 0)
	for k, v := range conf.Idx.Sites {
		for kk, vv := range v.Pages {
			//For each term, we get the number and presence
			//of the word for a particular page. The number
			//is currently discarded, because we can't rank
			//the relevance of pages.
			for i := range terms {
				_, isPresent := vv.WordCount[terms[i]]
				if isPresent {
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
