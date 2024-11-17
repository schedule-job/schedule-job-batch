package schedule

import (
	"errors"
	"strconv"
	"time"

	"github.com/schedule-job/schedule-job-batch/external/tool"
)

type everyWeeks struct {
	Weekdays []int `json:"weekdays"`
	Hour     int   `json:"hour"`
	Minute   int   `json:"minute"`
}

func everyWeeksToParams(params map[string]string) (*everyWeeks, error) {
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

	return &everyWeeks{
		Weekdays: weekdays,
		Hour:     hour,
		Minute:   minute,
	}, nil
}

func EveryWeeks(pivotTime time.Time, params map[string]string, _ interface{}) (*time.Time, error) {
	var payload, err = everyWeeksToParams(params)

	if err != nil {
		return nil, err
	}

	times := []time.Time{}
	current := int(pivotTime.Weekday())

	for _, week := range payload.Weekdays {
		day := week - current
		if day < 0 {
			day = day + 7
		}

		newTime := tool.NewUTCDate(pivotTime.Year(), pivotTime.Month(), pivotTime.Day()+day, payload.Hour, payload.Minute)

		if newTime.Before(pivotTime) {
			newTime = newTime.AddDate(0, 0, 7)
		}

		times = append(times, newTime)
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
