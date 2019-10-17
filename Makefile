export GO111MODULE := on
all: fmt lint vet build-dev test 
build-dev:
	go build ./cmd/graphite-writer-stats/graphite-writer-stats.go
fmt:
	go fmt ./...
vet:
	go vet ./...
lint:
	golint ./...
test:
	go test -v ./...
run:
	./graphite-writer-stats --brokers localhost:9092 --topic metrics --group graphite-writer-stats
docker-build:
	docker build . -t graphite-writer-stats -f build/Dockerfile
docker-kafka-start:
	docker-compose -f test/integration/docker-compose.yml up -d
docker-kafka-stop:
	docker-compose -f test/integration/docker-compose.yml down
docker-kafka-create-topic:
	docker exec kafka kafka-topics --create --zookeeper zookeeper:2181 --replication-factor 1 --partitions 5 --topic metrics
docker-start:
	docker run --rm --name graphite-writer-stats --network integration_kafka -p 8080:8080 -v $(shell pwd)/configs/:/app/configs/ graphite-writer-stats --brokers kafka:29092 --topic metrics -group graphite-writer-stats
docker-stop:
	docker stop graphite-writer-stats
docker-inject:
	docker exec -i kafka kafka-console-producer --broker-list localhost:9092 --topic metrics
