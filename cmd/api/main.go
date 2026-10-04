package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/abhishek8841/go-monolith/internal/config"
	"github.com/abhishek8841/go-monolith/internal/db"
	"github.com/abhishek8841/go-monolith/internal/handlers"
	"github.com/abhishek8841/go-monolith/internal/middleware"
)

func main() {
	cfg := config.MustLoad()

	db, err := db.Connect(cfg.DatabaseUrl)
	if err != nil {
		log.Fatalf("main.db.connect: %v", err)
	}

	// default is LevelInfo
	// the level which is selected that and all the levels with severity more that that will get logged
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	fmt.Println("Starting the server now...")

	// if we did http.HandleFunc then automatically the global mux would be registered which is considered bad practice so we create new mux and use that...
	mux := http.NewServeMux()
	lh := handlers.NewListingHandler(db, logger)
	mux.HandleFunc("GET /healthz", handlers.Health)
	mux.HandleFunc("GET /listings", lh.List)
	mux.HandleFunc("DELETE /listings/{id}", lh.Delete)
	mux.HandleFunc("POST /listings", lh.Create)

	handler := middleware.RequestId(mux)

	srv := http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      handler,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 30,
		IdleTimeout:  time.Second * 60,
	}

	// mux is also a type of handler that can store a lot of other handlers
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Server failed : %v", err)
	}
}
