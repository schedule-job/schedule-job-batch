package rule_based_replace

import (
	"strings"

	"github.com/schedule-job/schedule-job-batch/external/tool"
)

type ruleBasedReplaceOption struct {
	AllowRules []string
}

func getSupportFunc(text string) []string {
	return tool.FindWords(text, `(\[\:[^(]+\([^)(]*\)\:\])`)
}

func parseParams(params string) map[string]string {
	var words = tool.FindWords(params, `([^,=]+)=([^,=]+)`)
	var result = make(map[string]string)
	for i := 0; i < len(words); i += 2 {
		result[strings.TrimSpace(words[i])] = strings.TrimSpace(words[i+1])
	}
	return result
}

func parseNameWithParams(function string) (string, map[string]string) {
	var words = tool.FindWords(function, `\[\:([^(]+)\(([^)(]*)\)\:\]`)
	return words[0], parseParams(words[1])
}

type RuleBasedReplace struct {
	Options      ruleBasedReplaceOption
	rules        map[string]func(map[string]string, interface{}) (string, error)
	rulesOptions map[string]interface{}
}

func (t *RuleBasedReplace) AddRule(name string, f func(map[string]string, interface{}) (string, error), options interface{}) {
	if t.rules == nil {
		t.rules = make(map[string]func(map[string]string, interface{}) (string, error))
	}
	if t.rulesOptions == nil {
		t.rulesOptions = make(map[string]interface{})
	}
	t.rules[name] = f
	t.rulesOptions[name] = options
}

func (t *RuleBasedReplace) RuleBasedReplace(text string) string {
	var newText = text
	var words = getSupportFunc(text)
	for _, word := range words {
		name, params := parseNameWithParams(word)
		if t.rules[name] == nil {
			continue
		}
		if t.Options.AllowRules != nil && !tool.ContainsStringArray(t.Options.AllowRules, name) {
			continue
		}
		result, err := t.rules[name](params, t.rulesOptions[name])

		if err != nil {
			newText = strings.ReplaceAll(newText, word, "[:ERROR="+err.Error()+":]")
			continue
		}

		newText = strings.ReplaceAll(newText, word, result)
	}

	if newText != text {
		return t.RuleBasedReplace(newText)
	}
	return newText
}

var Replacer = RuleBasedReplace{}
