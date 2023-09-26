package opensearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/aaronland/go-pagination"
	"github.com/cenkalti/backoff/v4"
	go_opensearch "github.com/opensearch-project/opensearch-go/v2"
	go_opensearchtransport "github.com/opensearch-project/opensearch-go/v2/opensearchtransport"
	go_opensearchutil "github.com/opensearch-project/opensearch-go/v2/opensearchutil"
	"github.com/sfomuseum/go-libraryofcongress-database"
	"github.com/sfomuseum/go-timings"
)

type OpensearchV2Database struct {
	database.LibraryOfCongressDatabase
	client  *go_opensearch.Client
	index   string
	logger  *log.Logger
	workers int
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

	debug := false
	logger := log.New(io.Discard, "", 0)

	q := u.Query()

	os_endpoint := q.Get("endpoint")
	os_index := q.Get("index")
	str_workers := q.Get("workers")
	q_debug := q.Get("debug")

	if str_workers != "" {

		w, err := strconv.Atoi(str_workers)

		if err != nil {
			return nil, fmt.Errorf("Failed to parse workers, %w", err)
		}

		workers = w
	}

	if q_debug != "" {

		v, err := strconv.ParseBool(q_debug)

		if err != nil {
			return nil, fmt.Errorf("Failed to parse ?debug= parameter, %w", err)
		}

		debug = v
		logger = log.New(os.Stdout, "", 0)
	}

	retry := backoff.NewExponentialBackOff()

	os_cfg := go_opensearch.Config{
		Addresses: []string{
			os_endpoint,
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

	if debug {

		opensearch_logger := &go_opensearchtransport.TextLogger{
			Output:             os.Stdout,
			EnableRequestBody:  true,
			EnableResponseBody: true,
		}

		os_cfg.Logger = opensearch_logger
	}

	os_client, err := go_opensearch.NewClient(os_cfg)

	if err != nil {
		return nil, fmt.Errorf("Failed to create ES client, %w", err)
	}

	/*
		_, err = os_client.Indices.Create(os_index)

		if err != nil {
			return nil, fmt.Errorf("Failed to create index, %w", err)
		}
	*/

	opensearch_db := &OpensearchV2Database{
		client:  os_client,
		index:   os_index,
		workers: workers,
		logger:  logger,
	}

	return opensearch_db, nil
}

func (opensearch_db *OpensearchV2Database) Index(ctx context.Context, sources []*database.Source, monitor timings.Monitor) error {

	bi_cfg := go_opensearchutil.BulkIndexerConfig{
		Index:         opensearch_db.index,
		Client:        opensearch_db.client,
		NumWorkers:    opensearch_db.workers,
		FlushInterval: 30 * time.Second,
		OnError: func(ctx context.Context, err error) {
			opensearch_db.logger.Printf("OPENSEARCH bulk indexer reported an error: %v\n", err)
		},
		// OnFlushStart func(context.Context) context.Context // Called when the flush starts.
		OnFlushEnd: func(context.Context) {
			opensearch_db.logger.Printf("OPENSEARCH bulk indexer flush end")
		},
	}

	indexer, err := go_opensearchutil.NewBulkIndexer(bi_cfg)

	if err != nil {
		return fmt.Errorf("Failed to create bulk indexer, %w", err)
	}

	for _, src := range sources {

		err := opensearch_db.indexSource(ctx, indexer, src, monitor)

		if err != nil {
			return fmt.Errorf("Failed to index %s, %v", src.Label, err)
		}
	}

	err = indexer.Close(ctx)

	if err != nil {
		return fmt.Errorf("Failed to close indexer, %w", err)
	}

	stats := indexer.Stats()
	opensearch_db.logger.Printf("Stats %v\n", stats)

	return nil
}

func (opensearch_db *OpensearchV2Database) indexSource(ctx context.Context, indexer go_opensearchutil.BulkIndexer, src *database.Source, monitor timings.Monitor) error {

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
					opensearch_db.logger.Printf("ERROR: Failed to index %s, %s", doc_id, err)
				} else {
					opensearch_db.logger.Printf("ERROR: Failed to index %s, %s: %s", doc_id, res.Error.Type, res.Error.Reason)
				}
			},
		}

		err = indexer.Add(ctx, bulk_item)

		if err != nil {
			opensearch_db.logger.Printf("Failed to schedule %s, %v", doc_id, err)
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
