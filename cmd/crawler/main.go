package main

import (
	"net/http"
	"os"
	"time"

	"github.com/Martian-dev/skrull/internal/scheduler"
)

func main() {
	url := os.Args[1]

	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	m := map[string]struct{}{}
	scheduler := scheduler.Scheduler{
		HTTPClient: client,
		Visited:    m,
	}

	// TODO: prepare the parser
	// TODO: figure out how the results are stored (or atleast for now display the results)

	scheduler.Schedule(url)
}
