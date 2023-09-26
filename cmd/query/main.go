package main

import (
	"context"
	"log"

	_ "github.com/sfomuseum/go-libraryofcongress-database-opensearch"
	"github.com/sfomuseum/go-libraryofcongress-database/app/query"
)

func main() {

	ctx := context.Background()
	logger := log.Default()

	err := query.Run(ctx, logger)

	if err != nil {
		logger.Fatalf("Failed to run query, %v", err)
	}
}
