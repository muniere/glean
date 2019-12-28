package main

import (
	"log"

	"github.com/muniere/glean/internal/app/server"
)

func main() {
	srv := server.New("0.0.0.0", 2718)
	err := srv.Start()

	if err != nil {
		log.Fatal(err)
	}
}
