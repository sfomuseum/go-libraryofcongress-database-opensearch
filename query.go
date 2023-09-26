package opensearch

import (
	"github.com/sfomuseum/go-libraryofcongress-database"
)

type QueryResponse struct {
	Hits *QueryResponseHits `json:"hits"`
}

type QueryResponseHits struct {
	Total   *QueryResponseTotal   `json:"total"`
	Results []*database.QueryResult `json:"hits"`
}

type QueryResponseTotal struct {
	Value    int    `json:"value"`
	Relation string `json:"relation"`
}
