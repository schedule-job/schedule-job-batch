package schedule_test

import (
	"testing"
	"time"

	"github.com/schedule-job/schedule-job-batch/internal/schedule"
	"github.com/schedule-job/schedule-job-batch/internal/tool"
)

func TestEveryYear(t *testing.T) {
	schedule.Scheduler.AddSchedule("everyYear", schedule.EveryYear, nil)

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
				"months":   "[]",
				"weekdays": "[5]",
				"hour":     "9",
				"minute":   "0",
			},
			expected: tool.NewUTCDate(2022, 1, 1, 9, 0),
			hasError: true,
		},
		{
			name: "Test case 2",
			now:  tool.NewUTCDate(2022, 1, 1, 0, 0),
			payload: map[string]string{
				"months":   "[1]",
				"weekdays": "[]",
				"hour":     "9",
				"minute":   "0",
			},
			expected: tool.NewUTCDate(2022, 1, 1, 9, 0),
			hasError: true,
		},
		{
			name: "Test case 3",
			now:  tool.NewUTCDate(2022, 1, 1, 0, 0),
			payload: map[string]string{
				"months":   "[1]",
				"weekdays": "[1]",
				"hour":     "9",
				"minute":   "0",
			},
			expected: tool.NewUTCDate(2022, 1, 1, 9, 0),
			hasError: false,
		},
		{
			name: "Test case 4",
			now:  tool.NewUTCDate(2022, 1, 1, 9, 0),
			payload: map[string]string{
				"months":   "[1]",
				"weekdays": "[1]",
				"hour":     "9",
				"minute":   "0",
			},
			expected: tool.NewUTCDate(2022, 1, 1, 9, 0),
			hasError: false,
		},
		{
			name: "Test case 5",
			now:  tool.NewUTCDate(2022, 1, 1, 10, 0),
			payload: map[string]string{
				"months":   "[1]",
				"weekdays": "[1]",
				"hour":     "9",
				"minute":   "0",
			},
			expected: tool.NewUTCDate(2023, 1, 1, 9, 0),
			hasError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output, err := schedule.Scheduler.Schedule(tc.now, "everyYear", tc.payload)

			if (err != nil) != tc.hasError {
				t.Errorf("%v error = %v, wantErr %v", "everyYear", err, tc.hasError)
				return
			}
			if output != nil && !output.Equal(tc.expected) {
				t.Errorf("%v = %v, want %v", "everyYear", output, tc.expected)
			}
		})
	}
}
