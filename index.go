package distru

type Index struct {
	Pages []site //list of indexed webpages
}

type site struct {
	URL   string     //the link that identifies this Block
	Pages []sitePage //nonordered list of pages and their data on the server
}

type sitePage struct {
	Path    string   //path to page on the webserver (relative to root page)
	Links   []string //list of hyperlinks on the page
	Content string   //the content, temporarily replacing word lists
}

//sitePage constructor, which scrapes a webpage
func newSitePage(url string) *sitePage {
	body := fetch(url)      //get the body of the webpage
	links := getLinks(body) //get the links, as well

	//the wordlist should be added here, but that function doesn't exist yet
	//TODO

	page := sitePage{
		Path:    url,
		Links:   links,
		Content: body,
	}

	return &page
}
