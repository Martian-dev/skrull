// Package scheduler
package scheduler

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/Martian-dev/skrull/internal/fetcher"
)

type Scheduler struct {
	Visited    map[string]struct{}
	HTTPClient *http.Client
	mu         sync.Mutex
}

type URLJob struct {
	URL string
}

func (s *Scheduler) shouldVisit(url string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.Visited[url]; ok {
		return false
	}

	s.Visited[url] = struct{}{}
	return true
}

func worker(
	jobs chan URLJob,
	documents chan<- fetcher.Document,
	client *http.Client,
	tasks *sync.WaitGroup,
	s *Scheduler,
) {
	for job := range jobs {
		fmt.Println(job.URL) // just to print and show the url im getting
		doc, err := fetcher.FetchContent(job.URL, client)
		if err == nil {
			documents <- doc
			for _, link := range doc.Links {

				if !s.shouldVisit(link) {
					continue
				}

				tasks.Add(1)
				jobs <- URLJob{URL: link}
			}
		}
		tasks.Done()
	}
}

func (s *Scheduler) Schedule(url string) chan fetcher.Document {
	jobs := make(chan URLJob, 1000)
	documents := make(chan fetcher.Document, 1000)

	var tasks sync.WaitGroup

	// worker pool
	for range 10 {
		go worker(jobs, documents, s.HTTPClient, &tasks, s)
	}

	s.shouldVisit(url)
	tasks.Add(1)
	jobs <- URLJob{url}
	defer close(jobs)
	defer tasks.Wait()

	return documents
}
