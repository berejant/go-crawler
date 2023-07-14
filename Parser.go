package main

import (
	"errors"
	"github.com/spf13/afero"
	"golang.org/x/net/html"
	"io"
	"net/http"
	"strings"
)

// Parser - contains logic of handling single page: load, parse html, save to disk, any other extra
type Parser struct {

	// use abstract filesystem - to mock with Memory FS or change when it needs to Google cloud / AWS S3 / Cloudflare R2
	Fs afero.Fs
}

func (parser *Parser) Init() error {
	return parser.Fs.MkdirAll(".", 0755)
}

func (parser *Parser) ProcessPage(url string) (*ParserResponse, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("HTTP error " + resp.Status)
	}

	var fileWriter afero.File
	fileWriter, err = parser.Fs.Create(parser.makeFileName(url))
	if err != nil {
		return nil, err
	}
	defer fileWriter.Close()

	teeReader := io.TeeReader(resp.Body, fileWriter)

	var doc *html.Node

	// Use the html package to parse the response body from the request
	doc, err = html.Parse(teeReader)
	if err != nil {
		return nil, err
	}

	return &ParserResponse{
		Links: parser.extractLinks(doc),
	}, nil

}

func (parser *Parser) extractLinks(doc *html.Node) (links []string) {
	// Find and print all links on the web page
	var link func(*html.Node)
	link = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" && !strings.HasPrefix(a.Val, "#") {
					// adds a new link entry when the attribute matches
					links = append(links, a.Val)
				}
			}
		}

		// traverses the HTML of the webpage from the first child node
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			link(c)
		}
	}
	link(doc)

	return
}

func (parser *Parser) makeFileName(url string) string {
	url = strings.Replace(url, "https://", "", 1)
	url = strings.Replace(url, "http://", "", 1)
	url = strings.TrimSuffix(url, "/")
	url = strings.Replace(url, "/", "_", -1)

	if !strings.HasSuffix(url, ".html") {
		url += ".html"
	}
	return url
}
