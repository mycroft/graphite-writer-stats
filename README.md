# Graphite-Writer-Stats

This service reads a kafka topic for metrics in the plaintext graphite format to extract and count application names from them, and export these counters as a prometheus exporter.


## Building

Checkout the source code by cloning it or `go get` it.

```
go get github.com/criteo/graphite-writer-stats/...
$GOPATH/bin/graphite-writer-stats -h
```

If you clone the repository using `git clone`, you will be able to either `make build-dev` or `make docker-build`.


## Running

### Command line arguments

```
$ $GOPATH/bin/graphite-writer-stats -h
Usage of /home/mycroft/dev/go/bin/graphite-writer-stats:
  -brokers string
        Kafka bootstrap brokers to connect to, as a comma separated list (default "localhost:9092")
  -componentsNb uint
        number of components per extracted metric path. ex metric path: a.b.c.d with componentsNb=2 => a.b (default 3)
  -config string
        rule config path name (default "configs/rules.json")
  -endpoint string
        prometheus http endpoint name (default "/metrics")
  -group string
        Kafka consumer group id
  -oldest
        Kafka consumer consume initial offset from oldest (default true)
  -port uint
        prometheus http endpoint port (default 8080)
  -topic string
        Kafka topic to be consumed
```

### Rules configuration file

The rules configuration file is mandatory. It is used to find out which component(s) from the metrics' path will be used to count seen applications. The first rule matching will stop the processing.

There are 2 differents kind of rules: The path component matching's rules and the tag's rules.

#### Component matching rules

These rules will split the metric path in several components, will check that the different patterns matches with the extracted components (first pattern matches with first extracted component, second pattern with the second component, etc.). If so, the *applicationNamePosition*-nth component will be used as the application name.

Sample:

```
    {
      "name": "aggreg",
      "pattern": [
        "foo",
        "aggregated"
      ],
      "applicationNamePosition": 2
    }
```

#### Tag matching rules

These rules will check for existing graphite tags and will match if a tag with that name exists.

Sample:

```
    {
      "name": "by-tags",
      "use_tags": ["appname", "application"]
    }
```

### Example

```
$GOPATH/bin/graphite-writer-stats -brokers localhost:9092 -config config/rules.json -group writer-stats -topic metrics
```

You should be able to query from prometheus-format metrics on the http endpoint (default: 8080):

```
$ echo "foo.bar;appname=testaroo $RANDOM $(date +%s)" | make docker-inject

$ curl -s http://localhost:8080/metrics | grep foo.bar
metrics_path_total{application="testaroo",application_type="by-tags",metric_path="foo.bar"} 1
```

## Testing

Set up your kafka  with the docker compose: `make docker-kafka-start`.

You can inject some metrics (by default it inject to `metrics` topic):

```
echo "mymetric 42 000000000" | make docker-inject
```

You can read this topic by starting **graphite-writer-stats** by running it directly with `make run` or running it into its docker image by `make docker-start`.

Once done, stop everything by stopping kafka with `make docker-kafka-stop` and **graphite-writer-stats** with `make docker-stop`.
