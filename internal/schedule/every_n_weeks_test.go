package schedule_test

import (
	"testing"
	"time"

	"github.com/schedule-job/schedule-job-batch/internal/schedule"
	"github.com/schedule-job/schedule-job-batch/internal/tool"
)

func TestEveryNWeeks(t *testing.T) {
	schedule.Scheduler.AddSchedule("everyNWeeks", schedule.EveryNWeeks, nil)

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
				"startDate":     "2022-01-01-0-0",
				"intervalWeeks": "0",
			},
			expected: tool.NewUTCDate(2022, 1, 5, 9, 0),
			hasError: true,
		},
		{
			name: "Test case 2",
			now:  tool.NewUTCDate(2022, 1, 5, 10, 0),
			payload: map[string]string{
				"startDate":     "2022-01-01",
				"intervalWeeks": "1",
			},
			expected: tool.NewUTCDate(2022, 2, 5, 9, 0),
			hasError: true,
		},
		{
			name: "Test case 3",
			now:  tool.NewUTCDate(2022, 1, 10, 0, 0),
			payload: map[string]string{
				"startDate":     "2022-01-01-10-0",
				"intervalWeeks": "1",
			},
			expected: tool.NewUTCDate(2022, 1, 15, 10, 0),
			hasError: false,
		},
		{
			name: "Test case 4",
			now:  tool.NewUTCDate(2022, 1, 10, 0, 0),
			payload: map[string]string{
				"startDate":     "2022-01-01-10-0",
				"intervalWeeks": "2",
			},
			expected: tool.NewUTCDate(2022, 1, 15, 10, 0),
			hasError: false,
		},
		{
			name: "Test case 5",
			now:  tool.NewUTCDate(2022, 1, 15, 11, 0),
			payload: map[string]string{
				"startDate":     "2022-01-01-10-0",
				"intervalWeeks": "2",
			},
			expected: tool.NewUTCDate(2022, 1, 29, 10, 0),
			hasError: false,
		},
		{
			name: "Test case 6",
			now:  tool.NewUTCDate(2022, 1, 15, 10, 0),
			payload: map[string]string{
				"startDate":     "2022-01-01-10-0",
				"intervalWeeks": "2",
			},
			expected: tool.NewUTCDate(2022, 1, 15, 10, 0),
			hasError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output, err := schedule.Scheduler.Schedule(tc.now, "everyNWeeks", tc.payload)
			if (err != nil) != tc.hasError {
				t.Errorf("%v error = %v, wantErr %v", "everyNWeeks", err, tc.hasError)
				return
			}
			if output != nil && !output.Equal(tc.expected) {
				t.Errorf("%v = %v, want %v", "everyNWeeks", output, tc.expected)
			}
		})
	}
}
