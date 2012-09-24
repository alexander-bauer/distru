package main

import (
	"net/url"
)

type Index struct {
	Pages []Block //list of indexed webpages
	}

type Block struct {
	URL url //the link that identifies this Block
	Pages []Chunk //nonordered list of pages and their data on the server
	}

type Chunk struct {
	Path string //path to page on the webserver (relative to root page)
	Links []url //list of hyperlinks on the page
	Content string //the content, temporarily replacing word lists
	}
