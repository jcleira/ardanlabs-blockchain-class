SHELL := /bin/bash

# Bookeeping transactions
# curl -il -X GET http://localhost:8080/v1/sample
# curl -il -X GET http://localhost:9080/v1/node/sample
# curl -il -X GET http://localhost:8080/v1/tx/uncommitted/list
# curl -il -X GET http://localhost:8080/v1/start/mining

# ==============================================================================
# Local support

scratch:
	go run app/tooling/scratch/main.go

up:
	go run app/services/node/main.go -race | go run app/tooling/logfmt/main.go

up2:
	go run app/services/node/main.go -race --web-public-host 0.0.0.0:8280 --web-private-host 0.0.0.0:9280 | go run app/tooling/logfmt/main.go

down:
	kill -INT $(shell ps | grep "main -race" | grep -v grep | sed -n 1,1p | cut -c1-5)

down-ubuntu:
	kill -INT $(shell ps -x | grep "main -race" | sed -n 1,1p | cut -c3-7)

load:
	go run app/tooling/cli/main.go send -a kennedy -n 1 -t 0xbEE6ACE826eC3DE1B6349888B9151B92522F7F76 -v 100
	go run app/tooling/cli/main.go send -a pavel -n 1 -t 0xbEE6ACE826eC3DE1B6349888B9151B92522F7F76 -v 75
	go run app/tooling/cli/main.go send -a kennedy -n 2 -t 0x6Fe6CF3c8fF57c58d24BfC869668F48BCbDb3BD9 -v 150
	go run app/tooling/cli/main.go send -a pavel -n 2 -t 0xa988b1866EaBF72B4c53b592c97aAD8e4b9bDCC0 -v 125
	go run app/tooling/cli/main.go send -a kennedy -n 3 -t 0xa988b1866EaBF72B4c53b592c97aAD8e4b9bDCC0 -v 200
	go run app/tooling/cli/main.go send -a pavel -n 3 -t 0x6Fe6CF3c8fF57c58d24BfC869668F48BCbDb3BD9 -v 250


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