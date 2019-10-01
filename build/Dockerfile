FROM golang:1.13 as builder
WORKDIR /workspace
COPY go.sum go.mod ./
RUN go mod download
COPY . .
RUN go build ./cmd/graphite-writer-stats/graphite-writer-stats.go

FROM golang:1.13
WORKDIR /app
VOLUME configs/
COPY --from=builder /workspace/graphite-writer-stats .
ENTRYPOINT ["./graphite-writer-stats"]
CMD ["--brokers","kafka:9092","--topic","metrics","--group","graphite-writer"]
