package main

import (
	"sort"
	"time"
)

const (
	matchThreshold = 0.75 //Amount that a word must resemble a term to be a match
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
func (conf *config) Search(terms []string) (results []*page) {
	index := conf.Idx

	//Filter duplicate results through the use of maps.
	termsMap := make(map[string]*struct{}, len(terms))
	for i := range terms {
		termsMap[terms[i]] = nil
	}

	wordScore := uint64(0xffff / len(termsMap))

	if len(conf.Resources) > 0 {
		//Request indexes from all resources,
		//and trust their results. This is
		//done concurrently.
		workChan := make(chan bool)
		workers := len(conf.Resources)
		for _, resource := range conf.Resources {
			go func(resource string, workChan chan<- bool) {
				index.MergeRemote(resource, true, conf.ResTimeout)
				workChan <- true
			}(resource, workChan)
		}
		for _ = range workChan {
			workers--
			if workers == 0 {
				close(workChan)
			}
		}
	}

	bareresults := make(map[string]*page)
	for _, v := range conf.Idx.Sites {
		for kk, vv := range v.Pages {
			//For each term, we get the number and presence
			//of the word for a particular page. The number
			//is currently discarded, because we can't rank
			//the relevance of pages.
			for term := range termsMap {

				//Determine the "amount" that the words
				//match. This is determined by comparing
				//letters and their position. The total
				//score is the number of matches divided
				//by the search term's length.
				for w, c := range vv.WordCount {
					var resemblance int
					var longerTerm string
					var shorterTerm string
					if len(term) > len(w) {
						longerTerm = term
						shorterTerm = w
					} else {
						longerTerm = w
						shorterTerm = term
					}
					//We must iterate over the shorter
					//term, to avoid runtime errors.
					for i := range shorterTerm {
						if shorterTerm[i] == longerTerm[i] {
							resemblance++
						}
					}
					resemblanceM := float32(resemblance) / float32(len(longerTerm))
					if resemblanceM >= matchThreshold {
						if _, isFoundAlready := bareresults[kk]; !isFoundAlready {
							//If not already there, add it.

							//The multiplier will be:
							//12 divided by the number of hours since
							//the result was last used, capped at one
							//This means that it will be at most 1,
							//but only if the result is newer than
							//twelve hours old.
							timeM := 12 / time.Since(time.Unix(vv.Time, 0)).Hours()
							if timeM > 1 {
								timeM = 1
							}
							if timeM > 0 {
								//If the multiplier is negative, then the
								//timestamp is broken.
								vv.relevance = uint64(float64(wordScore) * timeM)
							}
							//Factor in the resemblance
							vv.relevance *= uint64(100 * resemblanceM)

							//And stamp it with the current time.
							vv.Time = time.Now().Unix()

							bareresults[kk] = vv
						} else {
							bareresults[kk].relevance += wordScore * uint64(c)
						}
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
	return c.Pages
}
