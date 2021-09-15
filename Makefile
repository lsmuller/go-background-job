include mocks.mk

TESTABLE_PACKAGES = `go list ./... | egrep -v 'mocks' | grep 'go-background-job/'`

setup:
	@go get github.com/golang/mock/gomock
	@go get github.com/golang/mock/mockgen

deps:
	@docker-compose --file boot/docker-compose.yaml --project-name=bitcoin-price-index-fetcher up --no-recreate -d redis

unit:
	@echo "Starting Unit Tests"
	@go test ${TESTABLE_PACKAGES} -tags=unit -coverprofile unit.coverprofile
	@echo "Finished Unit Tests"

run:
	@go run main.go

kill-deps:
	@docker-compose --project-name=bitcoin-price-index-fetcher -f ./boot/docker-compose.yaml down

mocks: all-mocks