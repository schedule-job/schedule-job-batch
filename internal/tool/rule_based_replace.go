package tool

func RuleBasedReplace(text string) string {
	return RuleBasedReplaceWithOption(text, nil)
}

func RuleBasedReplaceWithOption(text string, ignoreRule map[string]interface{}) string {
	return text
}
