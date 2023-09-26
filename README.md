# go-libraryofcongress-database-opensearch

Go package implementing the `sfomuseum/go-libraryofcongress-database.LibraryOfCongressDatabase` interface for use with OpenSearch.

## Documentation

Documentation is incomplete.
 
## Database URIs

```
opensearchv2://?endpoint=${OPENSEARCH_ENDPOINT}&index=${OPENSEARCH_INDEX}
```

Valid query parameters are:

| Name | Value | Notes | Required | 
| --- | --- | --- | --- |
| endpoint | string | The URI of your OpenSearch endpoint. | yes |
| index | string | The name of your OpenSearch index. | yes |
| debug | boolean | Enable verbose reporting to STDOUT. | no |
| create-index | bool | Create a new index. As of this writing this flag does _not_ assign a schema (described below) for the index. | no |
| workers | int | The number of workers to use when indexing records. Default is 10. | no |
| query-by | string | A flag to indicate whether to query by label or fulltext. Valid options are: 'label', 'text'. Default is 'label'. | no |

## Tools

```
$> make cli
go build -mod vendor -ldflags="-s -w" -o bin/server cmd/server/main.go
go build -mod vendor -ldflags="-s -w" -o bin/query cmd/query/main.go
go build -mod vendor -ldflags="-s -w" -o bin/index cmd/index/main.go
```

### index

```
$> ./bin/index -h
  -database-uri string
    		A valid sfomuseum/go-libraryofcongress-database URI.
  -lcnaf-data string
    	      The path to your LCNAF CSV data. If '-' then data will be read from STDIN.
  -lcsh-data string
    	     The path to your LCSH CSV data. If '-' then data will be read from STDIN.
```

For example:

```
$> ./bin/index \
	-database-uri 'opensearchv2://?endpoint=${OPENSEARCH_ENDPOINT}&index=${OPENSEARCH_INDEX}' \
	-lcsh-data /usr/local/data/lcsh.csv.bz2

```

It is also possible to index data from `STDIN` by specifying the string "-" as the `-lcsh-data` or `-lcnaf-data` URI to read.

For example, this command will stream and parse the contents of `https://id.loc.gov/download/lcsh.both.ndjson.zip` (using the `parse-lcsh` tool in the [sfomuseum/go-libraryofcongress](https://github.com/sfomuseum/go-libraryofcongress#parse-lcsh) package) and index each subject header in an OpenSearch database.

```
$> ./parse-lcsh https://id.loc.gov/download/lcsh.both.ndjson.zip | \
	./index \
	-database-uri 'opensearchv2://?endpoint=${OPENSEARCH_ENDPOINT}&index=${OPENSEARCH_INDEX}' \
	-lcsh-data -
```

_See notes about schemas below._

### query

```
$> ./bin/query -h
  -cursor-pagination
	Signal that pagination is cursor-based rather than countable.
  -database-uri string
    		A valid sfomuseum/go-libraryofcongress-database URI.
```

#### Example

Querying for a specific label:

```
$> ./bin/query \
	-database-uri 'opensearchv2://?endpoint=${OPENSEARCH_ENDPOINT}&index=${OPENSEARCH_INDEX}' \
	'Douglas DC-7 (Transport plane)' 

lcsh:sh2004006716 Douglas DC-7 (Transport plane)
```

Querying for a partial phrase across all labels can be done by appending a `query-by=text` query parameter to the database URI.

```
$> ./bin/query \
	-database-uri 'opensearchv2://?endpoint=${OPENSEARCH_ENDPOINT}&index=${OPENSEARCH_INDEX}&query-by=text' \
	'Montreal' 

lcnaf:no2020118720 Montreal (Kissock)
lcnaf:nr95002202 Montreal, Mary
lcnaf:n2021011867 Montreal lady
lcnaf:n2003074609 Montreal (Ship)
lcnaf:nr98040855 Montreal Library
lcnaf:n82024136 Bank of Montreal
lcnaf:no99015587 Montreal History Group
lcnaf:no2008037168 Montreal Cotton Company
lcnaf:no2002048901 Montreal, Steven R.
lcnaf:no2019017431 Montreal (Disco group)
lcnaf:no2005094329 Montreal (Musical group)
lcnaf:n2014183349 Montreal Medico-Chirurgical Society
lcnaf:n2010024889 CBC Montreal Choir
...and so on
```

_See notes about schemas, particularly when using the `?query-by=text` flag, below._

## Schemas

All the libraries and tools in this package assume an OpenSearch schema matching the [mappings.libraryofcongress](https://github.com/sfomuseum/es-sfomuseum-schema/blob/main/schema/7.4/mappings.libraryofcongress.json) mappings defined in the [sfomuseum/es-sfomuseum-schema](https://github.com/sfomuseum/es-sfomuseum-schema) package. Eventually that schema may be moved in to this package.

## See also

* https://github.com/sfomuseum/go-libraryofcongress
* https://github.com/sfomuseum/go-libraryofcongress-database
* https://pkg.go.dev/github.com/opensearch-project/opensearch-go