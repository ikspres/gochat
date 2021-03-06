package main

import (
	"github.com/ikspres/gochat/client"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <client nickname>", os.Args[0])
	}

	cli := client.NewClient(":6666", os.Args[1])
	defer cli.Close()

	cli.Go()
}
