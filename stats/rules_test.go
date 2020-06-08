package stats

import (
	"reflect"
	"testing"
)

func TestGetRules(t *testing.T) {

	var jsonRules = []byte(`{
  "rules": [
    {
      "name": "tag-rule",
      "use_tags": [
        "appname"
      ],
      "pattern": [],
      "applicationNamePosition": 0
    },
    {
      "name": "aggreg",
      "use_tags": [],
      "pattern": [
        "foo",
        "aggreg"
      ],
      "applicationNamePosition": 2
    },
    {
      "name": "anotheraggr",
      "use_tags": [],
      "pattern": [
        "foo",
        "anotheraggr"
      ],
      "applicationNamePosition": 2
    },
    {
      "name": "aggreg-all",
      "use_tags": [],
      "pattern": [
        "foo",
        "aggreg-all"
      ],
      "applicationNamePosition": 2
    },
    {
      "name": "legacy-bar",
      "use_tags": [],
      "pattern": [
        "prometheus",
        "bar"
      ],
      "applicationNamePosition": 1
    },
    {
      "name": "start-by-foo",
      "use_tags": [],
      "pattern": [
        "foo"
      ],
      "applicationNamePosition": 1
    },
    {
      "name": "start-by-app",
      "use_tags": [],
      "applicationNamePosition": 0
    }
  ]
}`)
	tagRule := Rule{"tag-rule", []string{"appname"}, []string{}, 0}
	aggregRule := Rule{"aggreg", []string{}, []string{"foo", "aggreg"}, 2}
	anotheraggrRule := Rule{"anotheraggr", []string{}, []string{"foo", "anotheraggr"}, 2}
	aggregAllRule := Rule{"aggreg-all", []string{}, []string{"foo", "aggreg-all"}, 2}
	legacybarRule := Rule{"legacy-bar", []string{}, []string{"prometheus", "bar"}, 1}
	startWithCriteoRule := Rule{"start-by-foo", []string{}, []string{"foo"}, 1}
	startbyAppRule := Rule{"start-by-app", []string{}, nil, 0}
	rulesExpected := []Rule{tagRule, aggregRule, anotheraggrRule, aggregAllRule, legacybarRule, startWithCriteoRule, startbyAppRule}
	rules, err := GetRulesFromBytes(jsonRules)

	for i := range rules.Rules {
		if !reflect.DeepEqual(rules.Rules[i], rulesExpected[i]) {
			t.Errorf("Failed to compare rules:\nExp. %v\nGot: %v", rules.Rules[i], rulesExpected[i])
		}
	}

	if (!reflect.DeepEqual(rules.Rules, rulesExpected)) || err != nil {
		t.Errorf("fail to parse rules : expected: '%v' actual: '%v', err: '%v'", rulesExpected, rules.Rules, err)
	}
}
func TestCheckRules(t *testing.T) {
	// Name, UseTags, Pattern, ApplicationNamePosition
	aggregRule := Rule{"aggreg", []string{}, []string{"foo", "aggreg"}, 2}
	anotheraggrRule := Rule{"anotheraggr", []string{}, []string{"foo", "anotheraggr"}, 2}
	aggregAllRule := Rule{"aggreg-all", []string{}, []string{"foo", "aggreg-all"}, 2}
	legacybarRule := Rule{"legacy-bar", []string{}, []string{"prometheus", "bar"}, 1}
	startWithCriteoRule := Rule{"start-by-foo", []string{}, []string{"foo"}, 1}
	startbyAppRule := Rule{"start-by-app", []string{}, nil, 0}
	rules := Rules{Rules: []Rule{aggregRule, anotheraggrRule, aggregAllRule, legacybarRule, startWithCriteoRule, startbyAppRule}}
	err := CheckRules(rules)
	if err != nil {
		t.Errorf("should not get the error: `%v`", err)
	}
	startbyAppRule = Rule{"", []string{}, nil, 0}
	rules = Rules{Rules: []Rule{aggregRule, anotheraggrRule, aggregAllRule, legacybarRule, startWithCriteoRule, startbyAppRule}}
	err = CheckRules(rules)
	if err == nil {
		t.Errorf("the rule should have a name: `%v`", err)
	}
	err = CheckRules(Rules{nil})
	if err == nil {
		t.Error("having a least one rule is mandatory")
	}
}

func TestErrorRules(t *testing.T) {
	rule := Rule{"rule1", []string{"foo"}, []string{"foo"}, 1}

	rules := Rules{Rules: []Rule{rule}}
	err := CheckRules(rules)
	if err == nil {
		t.Errorf("an error should happen when parsing rules: `%v`", rules)
	}
}
