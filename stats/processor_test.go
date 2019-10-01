package stats

import (
	"go.uber.org/zap/zaptest"
	"testing"
)

func TestProcess(t *testing.T) {
	logger := zaptest.NewLogger(t)
	aggregRule := Rule{"aggreg", []string{"foo", "aggreg"}, 2}
	anotheraggrRule := Rule{"anotheraggr", []string{"foo", "anotheraggr"}, 2}
	aggregAllRule := Rule{"aggreg-all", []string{"foo", "aggreg-all"}, 2}
	legacybarRule := Rule{"legacy-bar", []string{"prometheus", "bar"}, 1}
	startWithCriteoRule := Rule{"start-by-foo", []string{"foo"}, 1}
	startbyAppRule := Rule{"start-by-app", []string{}, 0}
	rulesTab := []Rule{aggregRule, anotheraggrRule, aggregAllRule, legacybarRule, startWithCriteoRule, startbyAppRule}

	stats := Stats{Logger: logger, MetricMetadata: MetricMetadata{
		Rules:        Rules{Rules: rulesTab},
		ComponentsNb: 3,
	}}

	metricDatapoint := []byte("foo.aggreg.cas.value 3.2 1498887")
	success := stats.Process(metricDatapoint)
	if !success {
		t.Errorf("failed to process '%v'", string(metricDatapoint))
	}
	metricDatapoint = nil
	success = stats.Process(metricDatapoint)
	if success {
		t.Errorf("process a nil metric should return false")
	}
	metricDatapoint = []byte("foo.498887")
	success = stats.Process(metricDatapoint)
	if success {
		t.Errorf("process a mal formed metric should return false '%v'", string(metricDatapoint))
	}
}
