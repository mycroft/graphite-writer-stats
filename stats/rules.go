package stats

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
)

type Rules struct {
	Rules []Rule `json:"rules"`
}
type Rule struct {
	Name                    string   `json:"name"`
	Pattern                 []string `json:"pattern"`
	ApplicationNamePosition uint     `json:"applicationNamePosition"`
}

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
