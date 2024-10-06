package schedule_test

import (
	"testing"
	"time"

	"github.com/schedule-job/schedule-job-batch/internal/schedule"
	"github.com/schedule-job/schedule-job-batch/internal/tool"
)

func TestEveryWeeks(t *testing.T) {
	schedule.Scheduler.AddSchedule("everyWeeks", schedule.EveryWeeks, nil)

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
				"weekdays": "[]",
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
				"weekdays": "[1]",
				"hour":     "9",
				"minute":   "0",
			},
			expected: tool.NewUTCDate(2022, 1, 3, 9, 0),
			hasError: false,
		},
		{
			name: "Test case 3",
			now:  tool.NewUTCDate(2022, 1, 1, 9, 0),
			payload: map[string]string{
				"weekdays": "[6]",
				"hour":     "9",
				"minute":   "0",
			},
			expected: tool.NewUTCDate(2022, 1, 1, 9, 0),
			hasError: false,
		},
		{
			name: "Test case 4",
			now:  tool.NewUTCDate(2022, 1, 1, 10, 0),
			payload: map[string]string{
				"weekdays": "[6]",
				"hour":     "9",
				"minute":   "0",
			},
			expected: tool.NewUTCDate(2022, 1, 8, 9, 0),
			hasError: false,
		},
		{
			name: "Test case 5",
			now:  tool.NewUTCDate(2022, 1, 1, 10, 0),
			payload: map[string]string{
				"weekdays": "[5,6]",
				"hour":     "9",
				"minute":   "0",
			},
			expected: tool.NewUTCDate(2022, 1, 7, 9, 0),
			hasError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output, err := schedule.Scheduler.Schedule(tc.now, "everyWeeks", tc.payload)

			if (err != nil) != tc.hasError {
				t.Errorf("%v error = %v, wantErr %v", "everyWeeks", err, tc.hasError)
				return
			}
			if output != nil && !output.Equal(tc.expected) {
				t.Errorf("%v = %v, want %v", "everyWeeks", output, tc.expected)
			}
		})
	}
}
