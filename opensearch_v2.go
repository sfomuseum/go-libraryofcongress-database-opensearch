package opensearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"time"

	"github.com/aaronland/go-pagination"
	"github.com/cenkalti/backoff/v4"
	go_opensearch "github.com/opensearch-project/opensearch-go/v2"
	go_opensearchutil "github.com/opensearch-project/opensearch-go/v2/opensearchutil"
	"github.com/sfomuseum/go-libraryofcongress-database"
	"github.com/sfomuseum/go-timings"
)

type OpensearchV2Database struct {
	database.LibraryOfCongressDatabase
	indexer go_opensearchutil.BulkIndexer
}

func init() {
	ctx := context.Background()
	database.RegisterLibraryOfCongressDatabase(ctx, "opensearch", NewOpensearchV2Database)
	database.RegisterLibraryOfCongressDatabase(ctx, "opensearchv2", NewOpensearchV2Database)
}

func NewOpensearchV2Database(ctx context.Context, uri string) (database.LibraryOfCongressDatabase, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URI, %w", err)
	}

	workers := 10

	q := u.Query()

	es_endpoint := q.Get("endpoint")
	es_index := q.Get("index")
	str_workers := q.Get("workers")

	if str_workers != "" {

		w, err := strconv.Atoi(str_workers)

		if err != nil {
			return nil, fmt.Errorf("Failed to parse workers, %w", err)
		}

		workers = w
	}

	retry := backoff.NewExponentialBackOff()

	es_cfg := go_opensearch.Config{
		Addresses: []string{
			es_endpoint,
		},

		RetryOnStatus: []int{502, 503, 504, 429},
		RetryBackoff: func(i int) time.Duration {
			if i == 1 {
				retry.Reset()
			}
			return retry.NextBackOff()
		},
		MaxRetries: 5,
	}

	es_client, err := go_opensearch.NewClient(es_cfg)

	if err != nil {
		return nil, fmt.Errorf("Failed to create ES client, %w", err)
	}

	/*
	_, err = es_client.Indices.Create(es_index)

	if err != nil {
		return nil, fmt.Errorf("Failed to create index, %w", err)
	}
	*/

	bi_cfg := go_opensearchutil.BulkIndexerConfig{
		Index:         es_index,
		Client:        es_client,
		NumWorkers:    workers,
		FlushInterval: 30 * time.Second,
	}

	indexer, err := go_opensearchutil.NewBulkIndexer(bi_cfg)

	if err != nil {
		return nil, fmt.Errorf("Failed to create bulk indexer, %w", err)
	}

	opensearch_db := &OpensearchV2Database{
		indexer: indexer,
	}

	return opensearch_db, nil
}

func (opensearch_db *OpensearchV2Database) Index(ctx context.Context, sources []*database.Source, monitor timings.Monitor) error {

	for _, src := range sources {

		err := opensearch_db.indexSource(ctx, src, monitor)

		if err != nil {
			return fmt.Errorf("Failed to index %s, %v", src.Label, err)
		}
	}

	return nil
}

func (opensearch_db *OpensearchV2Database) indexSource(ctx context.Context, src *database.Source, monitor timings.Monitor) error {

	cb := func(ctx context.Context, row map[string]string) error {

		doc := map[string]string{
			"id":     row["id"],
			"label":  row["label"],
			"source": src.Label,
		}

		doc_id := row["id"]

		enc_doc, err := json.Marshal(doc)

		if err != nil {
			return fmt.Errorf("Failed to marshal %s, %v", doc_id, err)
		}

		// log.Println(string(enc_doc))
		// continue

		bulk_item := go_opensearchutil.BulkIndexerItem{
			Action:     "index",
			DocumentID: doc_id,
			Body:       bytes.NewReader(enc_doc),

			OnSuccess: func(ctx context.Context, item go_opensearchutil.BulkIndexerItem, res go_opensearchutil.BulkIndexerResponseItem) {
				// log.Printf("Indexed %s\n", path)
			},

			OnFailure: func(ctx context.Context, item go_opensearchutil.BulkIndexerItem, res go_opensearchutil.BulkIndexerResponseItem, err error) {
				if err != nil {
					log.Printf("ERROR: Failed to index %s, %s", doc_id, err)
				} else {
					log.Printf("ERROR: Failed to index %s, %s: %s", doc_id, res.Error.Type, res.Error.Reason)
				}
			},
		}

		err = opensearch_db.indexer.Add(ctx, bulk_item)

		if err != nil {
			log.Printf("Failed to schedule %s, %v", doc_id, err)
			return nil
		}

		go monitor.Signal(ctx)
		return nil
	}

	return src.Index(ctx, cb)
}

func (opensearch_db *OpensearchV2Database) Query(ctx context.Context, q string, pg_opts pagination.Options) ([]*database.QueryResult, pagination.Results, error) {

	return nil, nil, fmt.Errorf("Not implemented")
}
