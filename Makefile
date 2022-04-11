test:
	GO111MODULE=on CGO_ENABLED=0 go test -v `go list ./...`

