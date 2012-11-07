package main

import ()

//Conf.Search returns the total number of results, and a []*page containing at most maxResults number of results.
func (conf *config) Search(terms []string) []*page {
	index := conf.Idx

	//Request indexes from all resources,
	//and trust their results.
	for i := range conf.Resources {
		index.MergeRemote(conf.Resources[i], true)
	}

	bareresults := make([]*page, 0)
	for _, v := range conf.Idx.Sites {
		for _, vv := range v.Pages {
			//For each term, we get the number and presence
			//of the word for a particular page. The number
			//is currently discarded, because we can't rank
			//the relevance of pages.
			for i := range terms {
				_, isPresent := vv.WordCount[terms[i]]
				if isPresent {
					bareresults = append(bareresults, vv)
				}
			}
		}
	}
	return bareresults
}
