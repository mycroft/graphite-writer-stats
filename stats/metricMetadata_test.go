package stats

import (
	"reflect"
	"testing"
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
	rule := Rule{"aggreg", []string{"foo", "aggreg"}, 2}
	isMatchedRule := isMatchingRule([]string{"foo", "aggreg", "d"}, rule)
	if !isMatchedRule {
		t.Error("should match the rule but for now it doesn't ")
	}
	rule = Rule{"aggreg", []string{"foo", "anotheraggr"}, 2}
	isMatchedRule = isMatchingRule([]string{"foo", "aggreg", "d"}, rule)
	if isMatchedRule {
		t.Error("should not match the rule but for now it doesn't ")
	}
	var rule2 Rule
	isMatchedRule = isMatchingRule([]string{"foo", "aggreg", "d"}, rule2)
	if !isMatchedRule {
		t.Error("should  match as pattern len equals 0 ! ")
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
	aggregRule := Rule{"aggreg", []string{"foo", "aggreg"}, 2}
	anotheraggrRule := Rule{"anotheraggr", []string{"foo", "anotheraggr"}, 2}
	aggregAllRule := Rule{"aggreg-all", []string{"foo", "aggreg-all"}, 2}
	legacybarRule := Rule{"legacy-bar", []string{"prometheus", "bar"}, 1}
	startWithCriteoRule := Rule{"start-by-foo", []string{"foo"}, 1}
	startbyAppRule := Rule{"start-by-app", nil, 0}
	rules := Rules{Rules: []Rule{aggregRule, anotheraggrRule, aggregAllRule, legacybarRule, startWithCriteoRule, startbyAppRule}}
	rule := getRule([]string{"foo", "aggreg", "myapp"}, rules)
	if rule.Name != aggregRule.Name {
		t.Error("should be an aggreg metric")
	}
}
