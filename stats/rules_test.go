package stats

import (
	"go.uber.org/zap/zaptest"
	"reflect"
	"testing"
)

func TestGetRules(t *testing.T) {

	var jsonRules = []byte(`{
  "rules": [
    {
      "name": "aggreg",
      "pattern": [
        "foo",
        "aggreg"
      ],
      "applicationNamePosition": 2
    },
    {
      "name": "anotheraggr",
      "pattern": [
        "foo",
        "anotheraggr"
      ],
      "applicationNamePosition": 2
    },
    {
      "name": "aggreg-all",
      "pattern": [
        "foo",
        "aggreg-all"
      ],
      "applicationNamePosition": 2
    },
    {
      "name": "legacy-bar",
      "pattern": [
        "prometheus",
        "bar"
      ],
      "applicationNamePosition": 1
    },
    {
      "name": "start-by-foo",
      "pattern": [
        "foo"
      ],
      "applicationNamePosition": 1
    },
    {
      "name": "start-by-app",
      "applicationNamePosition": 0
    }
  ]
}`)
	logger := zaptest.NewLogger(t)
	aggregRule := Rule{"aggreg", []string{"foo", "aggreg"}, 2}
	anotheraggrRule := Rule{"anotheraggr", []string{"foo", "anotheraggr"}, 2}
	aggregAllRule := Rule{"aggreg-all", []string{"foo", "aggreg-all"}, 2}
	legacybarRule := Rule{"legacy-bar", []string{"prometheus", "bar"}, 1}
	startWithCriteoRule := Rule{"start-by-foo", []string{"foo"}, 1}
	startbyAppRule := Rule{"start-by-app", nil, 0}
	rulesExpected := []Rule{aggregRule, anotheraggrRule, aggregAllRule, legacybarRule, startWithCriteoRule, startbyAppRule}
	rules, err := GetRulesFromBytes(logger, jsonRules)
	if (!reflect.DeepEqual(rules.Rules, rulesExpected)) || err != nil {
		t.Errorf("fail to parse rules : expected: '%v' actual: '%v', err: '%v'", rulesExpected, rules.Rules, err)
	}
}
func TestCheckRules(t *testing.T) {
	logger := zaptest.NewLogger(t)
	aggregRule := Rule{"aggreg", []string{"foo", "aggreg"}, 2}
	anotheraggrRule := Rule{"anotheraggr", []string{"foo", "anotheraggr"}, 2}
	aggregAllRule := Rule{"aggreg-all", []string{"foo", "aggreg-all"}, 2}
	legacybarRule := Rule{"legacy-bar", []string{"prometheus", "bar"}, 1}
	startWithCriteoRule := Rule{"start-by-foo", []string{"foo"}, 1}
	startbyAppRule := Rule{"start-by-app", nil, 0}
	rules := Rules{Rules: []Rule{aggregRule, anotheraggrRule, aggregAllRule, legacybarRule, startWithCriteoRule, startbyAppRule}}
	err := checkRules(logger, rules)
	if err != nil {
		t.Errorf("should not get the error: `%v`", err)
	}
	startbyAppRule = Rule{"", nil, 0}
	rules = Rules{Rules: []Rule{aggregRule, anotheraggrRule, aggregAllRule, legacybarRule, startWithCriteoRule, startbyAppRule}}
	err = checkRules(logger, rules)
	if err == nil {
		t.Errorf("the rule should have a name: `%v`", err)
	}
	err = checkRules(logger, Rules{nil})
	if err != nil {
		t.Error("rules is not mandatory")
	}
}
