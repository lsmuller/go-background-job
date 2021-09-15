# worker mocks #################################
all-mocks:
	@mockgen -source=worker/redis/worker.go -package mocks -destination mocks/redis_worker.go
	@mockgen -source=restclient/http.go -package mocks -destination mocks/http_client.go