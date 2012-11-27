package main

import (
	"time"
)

//Conf.Search returns the total number of results, and a []*page containing at most maxResults number of results. It returns all of the terms searched on, (omitting duplicates.)
func (conf *config) Search(terms []string) (results []*page, filteredTerms []string) {
	index := conf.Idx

	//Filter duplicate results through the use of maps.
	termsMap := make(map[string]*struct{}, len(terms))
	for i := range terms {
		termsMap[terms[i]] = nil
	}
	filteredTerms = make([]string, 0, len(termsMap))
	for k := range termsMap {
		filteredTerms = append(filteredTerms, k)
	}

	//Request indexes from all resources,
	//and trust their results.
	for i := range conf.Resources {
		index.MergeRemote(conf.Resources[i], true, conf.ResTimeout)
	}

	bareresults := make(map[string]*page)
	for _, v := range conf.Idx.Sites {
		for kk, vv := range v.Pages {
			//For each term, we get the number and presence
			//of the word for a particular page. The number
			//is currently discarded, because we can't rank
			//the relevance of pages.
			for i := range filteredTerms {
				_, isPresent := vv.WordCount[filteredTerms[i]]
				if isPresent {
					//I'm not sure why we set the time
					//here. TODO
					vv.Time = time.Now().String()
					bareresults[kk] = vv
				}
			}
		}
	}
	//The results should be sorted by relevance here. TODO
	results = make([]*page, 0, len(bareresults))
	for _, v := range bareresults {
		//We may want to speed this up by eliminating
		//append(). TODO
		results = append(results, v)
	}

	conf.Idx.Cache = append(conf.Idx.Cache, results...)
	return
}
