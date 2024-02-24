package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

const port = "2112"

func main() {
	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		metrics := []string{"metric1", "metric2", "metric3"}
		for _, metric := range metrics {
			value := rand.Intn(101)
			fmt.Fprintf(w, "otus_important_metrics[%s] %d\n", metric, value)
		}
	})

	server := &http.Server{
		Addr:         ":"+port,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	fmt.Printf("Starting server at port %s...", port)
	if err := server.ListenAndServe(); err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}