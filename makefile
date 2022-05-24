SHELL := /bin/bash

# Bookeeping transactions
# curl -il -X GET http://localhost:8080/v1/sample
# curl -il -X GET http://localhost:9080/v1/node/sample

# ==============================================================================
# Local support

scratch:
	go run app/tooling/scratch/main.go -race

up:
	go run app/services/node/main.go -race | go run app/tooling/logfmt/main.go

up2:
	go run app/services/node/main.go -race --web-public-host 0.0.0.0:8280 --web-private-host 0.0.0.0:9280 | go run app/tooling/logfmt/main.go

down:
	kill -INT $(shell ps | grep "main -race" | grep -v grep | sed -n 1,1p | cut -c1-5)

down-ubuntu:
	kill -INT $(shell ps -x | grep "main -race" | sed -n 1,1p | cut -c3-7)


# ==============================================================================
# Modules support

deps-reset:
	git checkout -- go.mod
	go mod tidy
	go mod vendor

tidy:
	go mod tidy
	go mod vendor

deps-upgrade:
	# go get $(go list -f '{{if not (or .Main .Indirect)}}{{.Path}}{{end}}' -m all)
	go get -u -v ./...
	go mod tidy
	go mod vendor

# ==============================================================================
# Running tests within the local computer

test:
	go test ./... -count=1
	staticcheck -checks=all ./...
