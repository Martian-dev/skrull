package main

import (
	"fmt"
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

	docs := scheduler.Schedule(url)

	for doc := range docs {
		fmt.Println(doc.URL)
	}
}
