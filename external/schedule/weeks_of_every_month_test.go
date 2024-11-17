package schedule_test

import (
	"testing"
	"time"

	"github.com/schedule-job/schedule-job-batch/external/schedule"
	"github.com/schedule-job/schedule-job-batch/external/tool"
)

func TestToDateByWeeksOfEveryMonth(t *testing.T) {
	schedule.Scheduler.AddSchedule("weeksOfEveryMonth", schedule.WeeksOfEveryMonth, nil)

	testCases := []struct {
		name     string
		now      time.Time
		payload  map[string]string
		expected time.Time
		hasError bool
	}{
		{
			name: "Test case 1",
			now:  tool.NewUTCDate(2022, 1, 1, 0, 0),
			payload: map[string]string{
				"monthlyWeekNumbers": "[]",
				"weekdays":           "[5]",
				"hour":               "9",
				"minute":             "0",
			},
			expected: tool.NewUTCDate(2022, 1, 1, 9, 0),
			hasError: true,
		},
		{
			name: "Test case 2",
			now:  tool.NewUTCDate(2022, 1, 1, 0, 0),
			payload: map[string]string{
				"monthlyWeekNumbers": "[1]",
				"weekdays":           "[1, 2, 3, 4]",
				"hour":               "9",
				"minute":             "0",
			},
			expected: tool.NewUTCDate(2022, 2, 1, 9, 0),
			hasError: false,
		},
		{
			name: "Test case 3",
			now:  tool.NewUTCDate(2022, 1, 1, 0, 0),
			payload: map[string]string{
				"monthlyWeekNumbers": "[1]",
				"weekdays":           "[6]",
				"hour":               "9",
				"minute":             "0",
			},
			expected: tool.NewUTCDate(2022, 1, 1, 9, 0),
			hasError: false,
		},
		{
			name: "Test case 4",
			now:  tool.NewUTCDate(2022, 1, 1, 0, 0),
			payload: map[string]string{
				"monthlyWeekNumbers": "[2]",
				"weekdays":           "[1]",
				"hour":               "9",
				"minute":             "0",
			},
			expected: tool.NewUTCDate(2022, 1, 3, 9, 0),
			hasError: false,
		},
		{
			name: "Test case 5",
			now:  tool.NewUTCDate(2022, 1, 1, 9, 0),
			payload: map[string]string{
				"monthlyWeekNumbers": "[1]",
				"weekdays":           "[6]",
				"hour":               "9",
				"minute":             "0",
			},
			expected: tool.NewUTCDate(2022, 1, 1, 9, 0),
			hasError: false,
		},
		{
			name: "Test case 6",
			now:  tool.NewUTCDate(2022, 1, 1, 9, 0),
			payload: map[string]string{
				"monthlyWeekNumbers": "[8]",
				"weekdays":           "[6]",
				"hour":               "9",
				"minute":             "0",
			},
			expected: tool.NewUTCDate(2022, 1, 1, 9, 0),
			hasError: true,
		},
		{
			name: "Test case 7",
			now:  tool.NewUTCDate(2022, 1, 1, 9, 0),
			payload: map[string]string{
				"monthlyWeekNumbers": "[1]",
				"weekdays":           "[1]",
				"hour":               "9",
				"minute":             "0",
			},
			expected: tool.NewUTCDate(2022, 1, 1, 9, 0),
			hasError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output, err := schedule.Scheduler.Schedule(tc.now, "weeksOfEveryMonth", tc.payload)
			if (err != nil) != tc.hasError {
				t.Errorf("%v error = %v, wantErr %v", "weeksOfEveryMonth", err, tc.hasError)
				return
			}
			if output != nil && !output.Equal(tc.expected) {
				t.Errorf("%v = %v, want %v", "weeksOfEveryMonth", output, tc.expected)
			}
		})
	}
}
