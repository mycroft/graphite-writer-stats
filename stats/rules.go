package stats

import (
	"encoding/json"
	"fmt"
)

// Rules is an array of Rule.
type Rules struct {
	Rules []Rule `json:"rules"`
}

// Rule structure
// Name will be used in prometheus metrics
// UseTags: If present and not empty, rule will match if any tags in list is present in metric.
// If not empty, pattern & applicationNamePosition will be ignored
// Pattern: Pattern to match the metric; If matching, the ApplicationNamePosition-nth will be used.
type Rule struct {
	Name                    string   `json:"name"`
	UseTags                 []string `json:"use_tags"`
	Pattern                 []string `json:"pattern"`
	ApplicationNamePosition uint     `json:"applicationNamePosition"`
}

// GetRulesFromBytes loads rules from json contents
func GetRulesFromBytes(jsonBytes []byte) (Rules, error) {
	var rules Rules

	err := json.Unmarshal(jsonBytes, &rules)
	if err != nil {
		return rules, err
	}

	err = CheckRules(rules)
	return rules, err
}

// CheckRules will return an error if no rule exists or invalid rule is found.
func CheckRules(rules Rules) error {
	if len(rules.Rules) <= 0 {
		return fmt.Errorf("no rules defined")
	}

	for i, rule := range rules.Rules {
		if len(rule.Name) <= 0 {
			return fmt.Errorf("Bad rule name `%v` at indice `%v`", rule.Name, i)
		}

		if len(rule.UseTags) > 0 && len(rule.Pattern) > 0 {
			return fmt.Errorf("rule `%v` `%v` has tags & patterns defined but are mutually exclusive", rule.Name, i)
		}
	}

	return nil
}
