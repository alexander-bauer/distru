package main

import (
	"sort"
	"time"
)

//resultContainer is used to contain an array of pages so that they can be sorted.
type resultContainer struct {
	Pages []*page
}

//Returns the length of c.Pages.
func (c *resultContainer) Len() int {
	return len(c.Pages)
}

//Returns true if the relevance of c.Pages[i] is less than or equal to that of c.Pages[j].
func (c *resultContainer) Less(i, j int) bool {
	return c.Pages[i].relevance <= c.Pages[j].relevance
}

//Swaps the indexes of i and j in c.Pages.
func (c *resultContainer) Swap(i, j int) {
	swap := c.Pages[i]
	c.Pages[i] = c.Pages[j]
	c.Pages[j] = swap
}

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

	wordScore := uint64(0xffff / len(filteredTerms))

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
			for _, term := range filteredTerms {
				instances, isPresent := vv.WordCount[term]
				if isPresent {
					//We set the time here so we
					//can put them in the cache
					//later.
					if _, isFoundAlready := bareresults[kk]; !isFoundAlready {
						//If not already there, add it.
						vv.Time = time.Now().String()
						bareresults[kk] = vv
					} else {
						bareresults[kk].relevance += wordScore * uint64(instances)
					}
				}
			}
		}
	}
	c := &resultContainer{
		Pages: make([]*page, 0, len(bareresults)),
	}
	for _, v := range bareresults {
		c.Pages = append(c.Pages, v)
	}

	//Sort c by relevance.
	sort.Sort(c)

	conf.Idx.Cache = append(conf.Idx.Cache, c.Pages...)
	return c.Pages, filteredTerms
}
