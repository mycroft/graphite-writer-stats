package stats

import (
	"testing"

	"github.com/Shopify/sarama"
	"go.uber.org/zap/zaptest"
)

func TestProcess(t *testing.T) {
	logger := zaptest.NewLogger(t)

	aggregRule := Rule{"aggreg", []string{}, []string{"foo", "aggreg"}, 2}
	anotheraggrRule := Rule{"anotheraggr", []string{}, []string{"foo", "anotheraggr"}, 2}
	aggregAllRule := Rule{"aggreg-all", []string{}, []string{"foo", "aggreg-all"}, 2}
	legacybarRule := Rule{"legacy-bar", []string{}, []string{"prometheus", "bar"}, 1}
	startWithCriteoRule := Rule{"start-by-foo", []string{}, []string{"foo"}, 1}
	startbyAppRule := Rule{"start-by-app", []string{}, []string{}, 0}
	rulesTab := []Rule{aggregRule, anotheraggrRule, aggregAllRule, legacybarRule, startWithCriteoRule, startbyAppRule}

	stats := Stats{MetricMetadata: MetricMetadata{
		Rules:        Rules{Rules: rulesTab},
		ComponentsNb: 3,
	}}

	metricDatapoint := []byte("foo.aggreg.cas.value 3.2 1498887")
	sm := new(sarama.ConsumerMessage)
	sm.Value = metricDatapoint
	err := stats.Process(logger, sm)
	if err != nil {
		t.Errorf("failed to process '%v'", string(metricDatapoint))
	}

	sm.Value = nil
	err = stats.Process(logger, sm)
	if err == nil {
		t.Errorf("process a nil metric should return an error.")
	}

	sm.Value = []byte("foo.498887")
	err = stats.Process(logger, sm)
	if err == nil {
		t.Errorf("process a malformed metric should return an error: '%v'.", string(sm.Value))
	}
}
