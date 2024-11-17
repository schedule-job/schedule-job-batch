package schedule

import (
	"errors"
	"strconv"
	"time"

	"github.com/schedule-job/schedule-job-batch/external/tool"
)

type everyYear struct {
	Months   []int `json:"months"`
	Weekdays []int `json:"weekdays"`
	Minute   int   `json:"minute"`
	Hour     int   `json:"hour"`
}

func everyYearToParams(params map[string]string) (*everyYear, error) {
	hour, hourErr := strconv.Atoi(params["hour"])
	if hourErr != nil {
		return nil, hourErr
	}

	minute, minuteErr := strconv.Atoi(params["minute"])
	if minuteErr != nil {
		return nil, minuteErr
	}

	if params["weekdays"] == "" {
		return nil, errors.New("Weekdays is empty")
	}

	weekdays := tool.ConvertToInArray(params["weekdays"])

	if params["months"] == "" {
		return nil, errors.New("Months is empty")
	}

	months := tool.ConvertToInArray(params["months"])

	return &everyYear{
		Months:   months,
		Weekdays: weekdays,
		Hour:     hour,
		Minute:   minute,
	}, nil
}

func EveryYear(pivotTime time.Time, params map[string]string, _ interface{}) (*time.Time, error) {
	var payload, err = everyYearToParams(params)

	if err != nil {
		return nil, err
	}

	times := []time.Time{}

	for _, month := range payload.Months {
		year := pivotTime.Year()
		if month < int(pivotTime.Month()) {
			year += 1
		}
		for _, day := range payload.Weekdays {
			nextDate := tool.NewUTCDate(year, time.Month(month), day, payload.Hour, payload.Minute)
			if nextDate.Before(pivotTime) && !pivotTime.Equal(nextDate) {
				nextDate = nextDate.AddDate(1, 0, 0)
			}

			times = append(times, nextDate)
		}
	}

	if len(times) == 0 {
		return nil, errors.New("Invalid Weekdays")
	}

	earliest := times[0]

	for _, date := range times[1:] {
		if date.Before(earliest) {
			earliest = date
		}
	}

	return &earliest, nil
}
