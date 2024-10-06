package schedule_test

import (
	"testing"
	"time"

	"github.com/schedule-job/schedule-job-batch/internal/schedule"
	"github.com/schedule-job/schedule-job-batch/internal/tool"
)

func TestDaily(t *testing.T) {
	schedule.Scheduler.AddSchedule("daily", schedule.Daily, nil)

	testCases := []struct {
		name     string
		now      time.Time
		payload  map[string]string
		expected time.Time
	}{
		{
			name: "Test case 1",
			now:  tool.NewUTCDate(2022, 1, 1, 0, 0),
			payload: map[string]string{
				"hour":   "9",
				"minute": "0",
			},
			expected: tool.NewUTCDate(2022, 1, 1, 9, 0),
		},
		{
			name: "Test case 2",
			now:  tool.NewUTCDate(2022, 1, 1, 10, 0),
			payload: map[string]string{
				"hour":   "9",
				"minute": "0",
			},
			expected: tool.NewUTCDate(2022, 1, 2, 9, 0),
		},
		{
			name: "Test case 3",
			now:  tool.NewUTCDate(2022, 1, 1, 10, 0),
			payload: map[string]string{
				"hour":   "11",
				"minute": "0",
			},
			expected: tool.NewUTCDate(2022, 1, 1, 11, 0),
		},
		{
			name: "Test case 4",
			now:  tool.NewUTCDate(2022, 1, 1, 10, 0),
			payload: map[string]string{
				"hour":   "10",
				"minute": "0",
			},
			expected: tool.NewUTCDate(2022, 1, 1, 10, 0),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			output, err := schedule.Scheduler.Schedule(tc.now, "daily", tc.payload)

			if err != nil {
				t.Errorf("%v error = %v", "daily", err)
				return
			}
			if output != nil && !output.Equal(tc.expected) {
				t.Errorf("%v = %v, want %v", "daily", output, tc.expected)
			}
		})
	}
}
