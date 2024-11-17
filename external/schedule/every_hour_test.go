package schedule_test

import (
	"testing"
	"time"

	"github.com/schedule-job/schedule-job-batch/external/schedule"
	"github.com/schedule-job/schedule-job-batch/external/tool"
)

func TestEveryHour(t *testing.T) {
	schedule.Scheduler.AddSchedule("everyHour", schedule.EveryHour, nil)

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
				"minute": "0",
			},
			expected: tool.NewUTCDate(2022, 1, 1, 0, 0),
		},
		{
			name: "Test case 2",
			now:  tool.NewUTCDate(2022, 1, 1, 10, 20),
			payload: map[string]string{
				"minute": "10",
			},
			expected: tool.NewUTCDate(2022, 1, 1, 11, 10),
		},
		{
			name: "Test case 3",
			now:  tool.NewUTCDate(2022, 1, 1, 0, 30),
			payload: map[string]string{
				"minute": "31",
			},
			expected: tool.NewUTCDate(2022, 1, 1, 0, 31),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			output, err := schedule.Scheduler.Schedule(tc.now, "everyHour", tc.payload)

			if err != nil {
				t.Errorf("%v error = %v", "everyHour", err)
				return
			}
			if output != nil && !output.Equal(tc.expected) {
				t.Errorf("%v = %v, want %v", "everyHour", output, tc.expected)
			}
		})
	}
}
