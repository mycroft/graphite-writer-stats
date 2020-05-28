package stats

import (
	"reflect"
	"testing"

	"go.uber.org/zap/zaptest"
)

func TestCheapEquals(t *testing.T) {
	isEquals := cheapEqual([]string{"a"}, []string{"a"})
	if !isEquals {
		t.Error("array should be equals ")
	}
	isEquals = cheapEqual([]string{"a"}, []string{"b"})
	if isEquals {
		t.Error("array should not be equals ")
	}
	isEquals = cheapEqual([]string{"a"}, nil)
	if isEquals {
		t.Error("array should not be equals ")
	}
	isEquals = cheapEqual(nil, []string{"a"})
	if isEquals {
		t.Error("array should not be equals ")
	}
	isEquals = cheapEqual(nil, nil)
	if !isEquals {
		t.Error("array should  be equals ")
	}
}
func TestIsMatchingRule(t *testing.T) {
	rule := Rule{"aggreg", []string{}, []string{"foo", "aggreg"}, 2}
	isMatchedRule := isMatchingRule([]string{"foo", "aggreg", "d"}, map[string]string{}, rule)
	if !isMatchedRule {
		t.Error("should match the rule but for now it doesn't ")
	}
	rule = Rule{"aggreg", []string{}, []string{"foo", "anotheraggr"}, 2}
	isMatchedRule = isMatchingRule([]string{"foo", "aggreg", "d"}, map[string]string{}, rule)
	if isMatchedRule {
		t.Error("should not match the rule but for now it doesn't ")
	}
	var rule2 Rule
	isMatchedRule = isMatchingRule([]string{"foo", "aggreg", "d"}, map[string]string{}, rule2)
	if !isMatchedRule {
		t.Error("should match as pattern len equals 0 ! ")
	}

	ruleTags := Rule{"tag", []string{"appname"}, []string{}, 1}
	isMatchedRule = isMatchingRule([]string{"foo", "aggreg", "d"}, map[string]string{"foo": "bar"}, ruleTags)
	if isMatchedRule {
		t.Error("rule tags should not match")
	}

	isMatchedRule = isMatchingRule([]string{"foo", "aggreg", "d"}, map[string]string{"appname": "foobar"}, ruleTags)
	if !isMatchedRule {
		t.Error("rule tags should match")
	}

}

func TestGetComponents(t *testing.T) {
	metricPath := "a.b.c.d"
	components := getComponents("a.b.c.d", 3)
	componentsExpected := []string{"a", "b", "c"}
	if !reflect.DeepEqual(components, componentsExpected) {
		t.Errorf("components not good for `%v` actual: `%v` expected:`%v`", metricPath, components, componentsExpected)
	}
	components = getComponents("a.b.c.d", 4)
	componentsExpected = []string{"a", "b", "c", "d"}
	if !reflect.DeepEqual(components, componentsExpected) {
		t.Errorf("components not good for `%v` actual: `%v` expected:`%v`", metricPath, components, componentsExpected)
	}
	//More than possible return the full components
	components = getComponents("a.b.c.d", 5)
	componentsExpected = []string{"a", "b", "c", "d"}
	if !reflect.DeepEqual(components, componentsExpected) {
		t.Errorf("components not good for `%v` actual: `%v` expected:`%v`", metricPath, components, componentsExpected)
	}
	//Less than possible  return empty array
	components = getComponents("a.b.c.d", 0)
	if len(components) != 0 {
		t.Errorf("components not good for `%v` actual: `%v` expected len of 0", metricPath, components)
	}
	//More than possible return the full components
	components = getComponents("a", 1)
	componentsExpected = []string{"a"}
	if !reflect.DeepEqual(components, componentsExpected) {
		t.Errorf("components not good for `%v` actual: `%v` expected:`%v`", metricPath, components, componentsExpected)
	}
}

func TestGetRule(t *testing.T) {
	aggregRule := Rule{"aggreg", []string{}, []string{"foo", "aggreg"}, 2}
	anotheraggrRule := Rule{"anotheraggr", []string{}, []string{"foo", "anotheraggr"}, 2}
	aggregAllRule := Rule{"aggreg-all", []string{}, []string{"foo", "aggreg-all"}, 2}
	legacybarRule := Rule{"legacy-bar", []string{}, []string{"prometheus", "bar"}, 1}
	startWithCriteoRule := Rule{"start-by-foo", []string{}, []string{"foo"}, 1}
	startbyAppRule := Rule{"start-by-app", []string{}, nil, 0}
	rules := Rules{Rules: []Rule{aggregRule, anotheraggrRule, aggregAllRule, legacybarRule, startWithCriteoRule, startbyAppRule}}
	rule := getRule([]string{"foo", "aggreg", "myapp"}, map[string]string{}, rules)
	if rule.Name != aggregRule.Name {
		t.Error("should be an aggreg metric")
	}
}

func TestLegacyAndTags(t *testing.T) {
	logger := zaptest.NewLogger(t)

	metric := Metric{Path: "a.b.c.d", Tags: map[string]string{"foo": "bar", "appname": "testaroo"}}
	rule := Rule{"rule-name", []string{"appname"}, []string{}, 0}
	rule1 := Rule{"rule-pb", []string{}, []string{"appname"}, 1}
	ruleLast := Rule{"last-rule", []string{}, []string{}, 0}

	stats := Stats{Logger: logger, MetricMetadata: MetricMetadata{
		Rules:        Rules{Rules: []Rule{rule, rule1, ruleLast}},
		ComponentsNb: 3,
	}}

	extractedMetric := stats.getMetric(metric.Path, metric.Tags)

	if extractedMetric.ApplicationName != "testaroo" {
		t.Errorf("Invalid appname found: %v != %v", extractedMetric.ApplicationName, "testaroo")
	}

	metric = Metric{Path: "a.b.c.d", Tags: map[string]string{"foo": "bar", "invalid": "testaroo"}}

	extractedMetric = stats.getMetric(metric.Path, metric.Tags)
	if extractedMetric.ApplicationName != "a" {
		t.Errorf("Invalid appname found: %v != %v", extractedMetric.ApplicationName, "a")
	}

	metric = Metric{Path: "appname.b.c.d", Tags: map[string]string{"foo": "bar", "invalid": "testaroo"}}

	extractedMetric = stats.getMetric(metric.Path, metric.Tags)
	if extractedMetric.ApplicationName != "b" {
		t.Errorf("Invalid appname found: %v != %v", extractedMetric.ApplicationName, "b")
	}

}
