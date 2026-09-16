package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/abhishek8841/go-monolith/internal/config"
	"github.com/abhishek8841/go-monolith/internal/handlers"
)

func main() {
	cfg := config.MustLoad()

	fmt.Println("Starting the server now...")

	// if we did http.HandleFunc then automatically the global mux would be registered which is considered bad practice so we create new mux and use that...
	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthz", handlers.Health)

	srv := http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      mux,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 30,
		IdleTimeout:  time.Second * 60,
	}

	// mux is also a type of handler that can store a lot of other handlers
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Server failed : %v", err)
	}
}
