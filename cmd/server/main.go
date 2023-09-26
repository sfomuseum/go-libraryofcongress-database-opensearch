package main

import (
	"context"
	"log"

	_ "github.com/sfomuseum/go-libraryofcongress-database-opensearch"
	"github.com/sfomuseum/go-libraryofcongress-database/app/server"
)

func main() {

	ctx := context.Background()
	logger := log.Default()

	err := server.Run(ctx, logger)

	if err != nil {
		logger.Fatalf("Failed to run server, %v", err)
	}
}
