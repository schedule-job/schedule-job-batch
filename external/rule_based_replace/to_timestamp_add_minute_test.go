package rule_based_replace_test

import (
	"testing"
	"time"

	"github.com/schedule-job/schedule-job-batch/external/rule_based_replace"
	"github.com/schedule-job/schedule-job-batch/external/tool"
)

func TestToTimestampAddMinute(t *testing.T) {
	pivot := tool.NewUTCDate(2021, time.January, 1, 0, 0)
	var option = rule_based_replace.ToTimestampOption{
		Pivot: &pivot,
	}
	rule_based_replace.Replacer.AddRule("toTimestampAddMinute", rule_based_replace.ToTimestampAddMinute, option)

	tests := []struct {
		name string
		text string
		want string
	}{
		{
			name: "test",
			text: "안녕하세요, [:toTimestampAddMinute(hour=4, minute=4, add=4):] 값으로 지정해주세요!",
			want: "안녕하세요, 1609474080000 값으로 지정해주세요!",
		},
		{
			name: "test",
			text: "안녕하세요, [:toTimestampAddMinute(hour=4, minute=4, add=-4):] 값으로 지정해주세요!",
			want: "안녕하세요, 1609473600000 값으로 지정해주세요!",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := rule_based_replace.Replacer.RuleBasedReplace(tt.text); got != tt.want {
				t.Errorf("toTimestampAddMinute() = %v, want %v", got, tt.want)
			}
		})
	}
}
