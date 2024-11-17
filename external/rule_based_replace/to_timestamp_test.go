package rule_based_replace_test

import (
	"testing"
	"time"

	"github.com/schedule-job/schedule-job-batch/external/rule_based_replace"
	"github.com/schedule-job/schedule-job-batch/external/tool"
)

func TestToTimestamp(t *testing.T) {
	pivot := tool.NewUTCDate(2021, time.January, 1, 0, 0)
	var option = rule_based_replace.ToTimestampOption{
		Pivot: &pivot,
	}
	rule_based_replace.Replacer.AddRule("toTimestamp", rule_based_replace.ToTimestamp, option)

	tests := []struct {
		name string
		text string
		want string
	}{
		{
			name: "test",
			text: "안녕하세요, [:toTimestamp(hour=4, minute=4):] 값으로 지정해주세요!",
			want: "안녕하세요, 1609473840000 값으로 지정해주세요!",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := rule_based_replace.Replacer.RuleBasedReplace(tt.text); got != tt.want {
				t.Errorf("toTimestamp() = %v, want %v", got, tt.want)
			}
		})
	}
}
