run:
	go run main.go rest

run-redis:
	docker compose up -d redis

test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

test-burst:
	chmod -R 775 ./test_rate_limiter.sh
	./test_rate_limiter.sh

run-with-docker:
	docker compose up 