package main

import (
	"fmt"
	"log"
	"os"

	"github.com/abhishek8841/go-monolith/internal/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
    _ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	fmt.Println("Running migration")
	// fmt.Println(os.Args)
	// from os.Args slice output the second ele to get the command passed in argument while running the run command ?
	if len(os.Args) < 2 {
		log.Fatal("usage: make migrate <up|down>")
	}
	cfg := config.MustLoad()

	m, err := migrate.New(
		"file://migrations",
		cfg.DatabaseUrl)
	if err != nil {
		log.Fatalf("migration.new: %v", err)
	}

	switch os.Args[1] {
	case "up":
		if err := m.Up(); err != nil {
			log.Fatal(err)
		}
		fmt.Println("up")
	case "down":
		if err := m.Down(); err != nil {
			log.Fatal(err)
		}
		fmt.Println("down")
	default:
		log.Fatalf("unknown command %s", os.Args[1])
	}
}
