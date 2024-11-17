package schedule

import (
	"errors"
	"strconv"
	"time"

	"github.com/schedule-job/schedule-job-batch/external/tool"
)

type everyMonth struct {
	Days   []int `json:"days"`
	Hour   int   `json:"hour"`
	Minute int   `json:"minute"`
}

func everyMonthToParams(params map[string]string) (*everyMonth, error) {
	hour, hourErr := strconv.Atoi(params["hour"])
	if hourErr != nil {
		return nil, hourErr
	}

	minute, minuteErr := strconv.Atoi(params["minute"])
	if minuteErr != nil {
		return nil, minuteErr
	}

	if params["days"] == "" {
		return nil, errors.New("Days is empty")
	}

	days := tool.ConvertToInArray(params["days"])

	return &everyMonth{
		Days:   days,
		Hour:   hour,
		Minute: minute,
	}, nil
}

func EveryMonth(pivotTime time.Time, params map[string]string, _ interface{}) (*time.Time, error) {
	var payload, err = everyMonthToParams(params)

	if err != nil {
		return nil, err
	}

	times := []time.Time{}

	for _, day := range payload.Days {
		newTime := tool.NewUTCDate(pivotTime.Year(), pivotTime.Month(), day, payload.Hour, payload.Minute)

		if newTime.Before(pivotTime) {
			continue
		}

		times = append(times, newTime)
	}

	if len(payload.Days) > 0 {
		newTime := tool.NewUTCDate(pivotTime.Year(), pivotTime.Month()+1, payload.Days[0], payload.Hour, payload.Minute)
		times = append(times, newTime)
	}

	if len(times) == 0 {
		return nil, errors.New("no valid date")
	}

	earliest := times[0]

	for _, date := range times[1:] {
		if date.Before(earliest) {
			earliest = date
		}
	}

	return &earliest, nil
}
