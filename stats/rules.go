package stats

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
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
func GetRulesFromBytes(logger *zap.Logger, jsonBytes []byte) (Rules, error) {
	var rules Rules
	err := json.Unmarshal(jsonBytes, &rules)
	err = checkRules(logger, rules)
	return rules, err
}

func checkRules(logger *zap.Logger, rules Rules) error {
	var err error
	if len(rules.Rules) <= 0 {
		zap.L().Warn("No rules defined.")
	}
	for i, rule := range rules.Rules {
		if len(rule.Name) <= 0 {
			err = fmt.Errorf("Bad rule name `%v` at indice `%v`", rule.Name, i)
		}
	}
	return err
}
