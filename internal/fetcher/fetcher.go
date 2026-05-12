// Package fetcher
package fetcher

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

type Document struct {
	URL     *url.URL
	Title   string
	Content string
	Links   []string
}

func FetchContent(url string, client *http.Client) (Document, error) {
	rsp, err := client.Get(url)
	if err != nil {
		fmt.Println("Error fetching URL:", err)
		log.Fatal(err)
		return Document{}, err
	}

	defer func() {
		_ = rsp.Body.Close()
	}()

	baseURL := rsp.Request.URL

	doc, err := html.Parse(rsp.Body)
	if err != nil {
		fmt.Println("Error reading body:", err)
		log.Fatal(err)
		return Document{}, err
	}

	urls := extractURLs(doc, baseURL)

	return Document{
		Links: urls,
		URL:   baseURL,
	}, nil
}

func extractURLs(n *html.Node, baseURL *url.URL) []string {
	var urls []string

	var walk func(*html.Node)
	walk = func(node *html.Node) {
		if node.Type == html.ElementNode {
			for _, attr := range node.Attr {
				// common URL attributes
				if isURLAttr(attr.Key) {
					link := strings.TrimSpace(attr.Val)

					if link == "" {
						continue
					}

					// resolve relative URLs
					parsed, err := url.Parse(link)
					if err == nil {
						link = baseURL.ResolveReference(parsed).String()
					}

					urls = append(urls, link)
				}
			}
		}

		for c := node.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}

	walk(n)

	return urls
}

func isURLAttr(attr string) bool {
	switch attr {
	case "href", "src", "action", "poster":
		return true
	default:
		return false
	}
}
