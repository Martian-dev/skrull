// Package scheduler
package scheduler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Martian-dev/skrull/internal/fetcher"
)

type Scheduler struct {
	Visited    map[string]struct{}
	HTTPClient *http.Client
}

type URLJob struct {
	URL string
}

func (s *Scheduler) dedupe(input []string) []string {
	var result []string

	for _, v := range input {
		if _, ok := s.Visited[v]; ok {
			continue
		}

		s.Visited[v] = struct{}{}
		result = append(result, v)
	}

	return result
}

func (s *Scheduler) Schedule(url string) {
	jobs := make(chan URLJob, 1000)
	documents := make(chan fetcher.Document, 1000)

	jobs <- URLJob{url}

	for {
		select {
		case job := <-jobs:
			// schedule the job and handle its results
			fmt.Println(job.URL) // just to print and show the url im getting
			doc, err := fetcher.FetchContent(job.URL, s.HTTPClient)
			if err != nil {
				fmt.Println("Error fetching content", err)
				log.Fatal(err)
			}
			documents <- doc

		case doc := <-documents:
			// process parsed document
			dedupedURL := s.dedupe(doc.Links)
			doc.Links = dedupedURL
			for _, link := range dedupedURL {
				jobs <- URLJob{link}
			}

		default:
			fmt.Println("Done")
		}
	}
}
