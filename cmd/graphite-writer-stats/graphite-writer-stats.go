package main

import (
	"flag"
	"github.com/criteo/graphite-writer-stats/input"
	"github.com/criteo/graphite-writer-stats/prometheus"
	"github.com/criteo/graphite-writer-stats/stats"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"strconv"
)

var (
	brokers      = flag.String("brokers", "localhost:9092", "Kafka bootstrap brokers to connect to, as a comma separated list")
	group        = flag.String("group", "", "Kafka consumer group id")
	topic        = flag.String("topic", "", "Kafka topic to be consumed")
	oldest       = flag.Bool("oldest", true, "Kafka consumer consume initial offset from oldest")
	componentsNb = flag.Uint("componentsNb", 3, "number of components per extracted metric path. ex metric path: a.b.c.d with componentsNb=2 => a.b")
	port         = flag.Uint("port", 8080, "prometheus http endpoint port")
	endpoint     = flag.String("endpoint", "/metrics", "prometheus http endpoint name")
	config       = flag.String("config", "configs/rules.json", "rule config path name")
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	flag.Parse()
	if len(*brokers) == 0 {
		logger.Fatal("no Kafka bootstrap brokers defined, please set the -brokers flag")
	}
	if len(*topic) == 0 {
		logger.Fatal("no Kafka topic given to be consumed, please set the -topic flag")
	}
	if len(*group) == 0 {
		logger.Fatal("no Kafka consumer group defined, please set the -group flag")
	}
	if *componentsNb <= 0 {
		logger.Fatal("ComponentsNb should be > 0")
	}
	jsonRules, err := ioutil.ReadFile(*config)
	rules, err := stats.GetRulesFromBytes(logger, jsonRules)
	if err != nil {
		logger.Fatal("bad config rule.", zap.String("configFile", *config), zap.Error(err))
	}
	stats := stats.Stats{Logger: logger, MetricMetadata: stats.MetricMetadata{ComponentsNb: *componentsNb, Rules: rules}}

	kafka := input.SetupConsumer(logger, *oldest, *group, *brokers, *topic, stats)
	kafka.Run()

	go func() {
		portBinding := ":" + strconv.Itoa(int(*port))
		http.Handle(*endpoint, prometheus.GetPrometheusHTTPHandler())
		http.Handle("/", kafka.GetStatusHTTPHandler())
		err := http.ListenAndServe(portBinding, nil)
		if err != nil {
			logger.Panic("could not set up HTTP server", zap.Error(err))
		}
	}()

	kafka.Wait()
	kafka.Close()
}