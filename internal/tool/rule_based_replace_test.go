package tool_test

import (
	"errors"
	"strconv"
	"testing"

	"github.com/schedule-job/schedule-job-batch/internal/tool"
)

func TestRuleBasedReplace(t *testing.T) {

	var newTool = tool.Tool{}

	newTool.AddRule("add", func(params map[string]string, options interface{}) (string, error) {
		if params["a"] == "" || params["b"] == "" {
			return "", errors.New("a or b is empty")
		}
		a, aErr := strconv.Atoi(params["a"])
		if aErr != nil {
			return "", aErr
		}
		b, bErr := strconv.Atoi(params["b"])
		if bErr != nil {
			return "", bErr
		}
		return strconv.Itoa(a + b), nil
	}, nil)

	tests := []struct {
		name string
		text string
		want string
	}{
		{
			name: "test",
			text: "안녕하세요 4 + 4는 [:add(a=4, b=4):] 입니다.",
			want: "안녕하세요 4 + 4는 8 입니다.",
		},
		{
			name: "test",
			text: "안녕하세요 4 + 4는 [:add(a=4):] 입니다.",
			want: "안녕하세요 4 + 4는 [:ERROR=a or b is empty:] 입니다.",
		},
		{
			name: "test",
			text: "안녕하세요 4 + 4는 [:add(a=4,b=[:add(a=4,b=8):]):] 입니다.",
			want: "안녕하세요 4 + 4는 16 입니다.",
		},
		{
			name: "test",
			text: "안녕하세요 4 + 4는 [:add(a=[:add(a=4,b=8):],b=[:add(a=4,b=8):]):] 입니다.",
			want: "안녕하세요 4 + 4는 24 입니다.",
		},
		{
			name: "test",
			text: "안녕하세요 4 + 4는 [:add(a=[:add(a=4,b=8):],b=[:add(a=4,b=[:add(a=[:add(a=4,b=8):],b=[:add(a=4,b=8):]):]):]):] 입니다.",
			want: "안녕하세요 4 + 4는 40 입니다.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newTool.RuleBasedReplace(tt.text); got != tt.want {
				t.Errorf("RuleBasedReplace() = %v, want %v", got, tt.want)
			}
		})
	}
}
