module github.com/sfomuseum/go-libraryofcongress-database-opensearch

go 1.21.1

// replace statement for opensearch-go necessary until this get	sorted...
// go: github.com/opensearch-project/opensearch-go/v2@v2.3.0 requires
// github.com/aws/aws-sdk-go-v2/config@v1.18.25 requires
// github.com/aws/aws-sdk-go-v2/feature/ec2/imds@v1.13.3: reading github.com/aws/aws-sdk-go-v2/feature/ec2/imds/go.mod at revision feature/ec2/imds/v1.13.3: unknown revision refs/tags/feature/ec2/imds/v1.13.3

replace github.com/opensearch-project/opensearch-go/v2 => github.com/opensearch-project/opensearch-go/v2 v2.2.1-0.20230919181059-011b99e67c6e

require (
	github.com/aaronland/go-pagination v0.2.0
	github.com/cenkalti/backoff/v4 v4.2.1
	github.com/opensearch-project/opensearch-go/v2 v2.0.0-00010101000000-000000000000
	github.com/sfomuseum/go-libraryofcongress-database v0.1.2
	github.com/sfomuseum/go-timings v1.2.1
)

require (
	github.com/aaronland/go-http-bootstrap v0.4.0 // indirect
	github.com/aaronland/go-http-rewrite v1.1.0 // indirect
	github.com/aaronland/go-http-sanitize v0.0.8 // indirect
	github.com/aaronland/go-http-server v1.2.0 // indirect
	github.com/aaronland/go-http-static v0.0.3 // indirect
	github.com/aaronland/go-log/v2 v2.0.0 // indirect
	github.com/aaronland/go-roster v1.0.0 // indirect
	github.com/akrylysov/algnhsa v1.0.0 // indirect
	github.com/aws/aws-lambda-go v1.37.0 // indirect
	github.com/jtacoma/uritemplates v1.0.0 // indirect
	github.com/sfomuseum/go-csvdict v1.0.0 // indirect
	github.com/sfomuseum/go-flags v0.10.0 // indirect
	github.com/sfomuseum/go-http-rollup v0.0.2 // indirect
	github.com/sfomuseum/iso8601duration v1.1.0 // indirect
	github.com/tdewolff/minify/v2 v2.12.4 // indirect
	github.com/tdewolff/parse/v2 v2.6.4 // indirect
	github.com/whosonfirst/go-sanitize v0.1.0 // indirect
	golang.org/x/net v0.13.0 // indirect
)
