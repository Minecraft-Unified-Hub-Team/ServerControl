proto:
	docker container run -v .:/mnt/grpc --rm noname0443/go_proto_builder /bin/sh -c "cd /mnt/grpc; protoc -I . -I /usr/lib/ --go_out=. --go-grpc_out=. ./internal/api/*.proto"

run:
	go run cmd/main.go

integration-tests: proto
	docker compose up --wait -d --build
	cd tests
	go test ./... -count=1 -v

unit-tests:
	go test internal/...

cleanup:
	docker compose down
