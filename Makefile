deps:
	@docker-compose --file boot/docker-compose.yaml --project-name=bitcoin-price-index-fetcher up --no-recreate -d redis

run:
	@go run main.go

kill-deps:
	@docker-compose --project-name=bitcoin-price-index-fetcher -f ./boot/docker-compose.yaml down