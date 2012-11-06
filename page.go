package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"strings"
)

//getPage is a complex constructor for the page object. It appends path to target in order to get the target webpage. It then uses http.Get to get the body of that webpage, which it then uses regexp to scrape for links. Those links are sorted into internal and external. The internal links are resolved to be absolute (internal) links on the webserver, and then returned, without duplicates, as a map[string]struct{}. All unique external links on the page are returned in the second map[string]struct{}.
func getPage(target, path string, client http.Client) (*page, map[string]struct{}, map[string]struct{}) {
	//Parse the target URI, return empty if it fails.
	accessURI, err := url.ParseRequestURI(target + path)
	if err != nil {
		return nil, nil, nil
	}

	//Get the content of the webpage via HTTP, using the
	//existing http.Client, and return blank if it fails.
	resp, err := http.Get(accessURI.String())
	if err != nil {
		return nil, nil, nil
	}
	defer resp.Body.Close()
	//Get the body of the request as a []byte.
	b, err := ioutil.ReadAll(resp.Body)
	//Convert to string real quick.
	body := string(b)

	//if there is no title, show the name of the url
	title := target
	title = title[len("http://"):]
	titlepattern, err := regexp.Compile("<title>.*</title>")
	if err != nil {
		return nil, nil, nil
	}
	//Find the leftmost title tag
	titleb := titlepattern.Find(b)
	//and cut out the html tags, if
	//title is present at all.
	if titleb != nil {
		title = string(titleb[len("<title>") : len(titleb)-len("</title>")])
	}

	//Now we're going to move on to parsing the links.
	pattern, err := regexp.Compile("href=['\"]?([^'\" >]+)")
	if err != nil {
		return nil, nil, nil
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
	p, err := regexp.Compile("<([^>]*)>|\n|\t|&[a-z]+|[.,]+ |;|\u0009")
	if err != nil {
		return nil, nil, nil
	}
	//Apply the pattern and split on spaces.
	content := bytes.Split(p.ReplaceAll(b, []byte("")), []byte(" "))
	wc := make(map[string]int)

	//For every word...
	for i := range content {
		word := string(content[i])
		//if the word is less than two characters long
		//or is one of the listed common words,
		if len(word) < 2 || word == "is" || word == "or" || word == "a" || word == "and" || word == "the" || word == "are" || word == "of" || word == "to" {
			//then skip it.
			continue
		}
		//Otheriwse, increment that word's counter by one.
		wc[word] += 1
	}

	return &page{
		Title:     string(title),
		Link:      target + path,
		WordCount: wc,
	}, internalLinks, externalLinks
}

func join(source, target string) string {
	if path.IsAbs(target) {
		return target
	}
	return path.Join(path.Dir(source), target)
}
