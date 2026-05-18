// Package fetcher
package fetcher

import (
	"fmt"
	"net/http"
	"net/url"
	"path"
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
		return Document{}, err
	}

	defer func() {
		_ = rsp.Body.Close()
	}()

	baseURL := rsp.Request.URL

	doc, err := html.Parse(rsp.Body)
	if err != nil {
		fmt.Println("Error reading body:", err)
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
	seen := make(map[string]struct{})

	rejectedExt := map[string]struct{}{
		".js":   {},
		".css":  {},
		".png":  {},
		".jpg":  {},
		".jpeg": {},
		".svg":  {},
	}

	var walk func(*html.Node)

	walk = func(node *html.Node) {
		if node.Type == html.ElementNode {
			for _, attr := range node.Attr {

				if !isURLAttr(attr.Key) {
					continue
				}

				link := strings.TrimSpace(attr.Val)

				if link == "" {
					continue
				}

				// skip fragments
				if strings.HasPrefix(link, "#") {
					continue
				}
				// reject obvious bad schemes early
				lower := strings.ToLower(link)

				if strings.HasPrefix(lower, "mailto:") ||
					strings.HasPrefix(lower, "javascript:") {
					continue
				}

				parsed, err := url.Parse(link)
				if err != nil {
					continue
				}

				// skip mailto:, javascript:, etc
				if parsed.Scheme != "" &&
					parsed.Scheme != "http" &&
					parsed.Scheme != "https" {
					continue
				}

				// convert relative -> absolute
				resolved := baseURL.ResolveReference(parsed)

				// only same host
				if resolved.Host != baseURL.Host {
					continue
				}

				// optional normalization
				resolved.Fragment = ""
				// reject static assets by extension
				ext := strings.ToLower(path.Ext(resolved.Path))

				if _, blocked := rejectedExt[ext]; blocked {
					continue
				}

				finalURL := resolved.String()

				// dedupe
				if _, exists := seen[finalURL]; exists {
					continue
				}

				seen[finalURL] = struct{}{}
				urls = append(urls, finalURL)
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
