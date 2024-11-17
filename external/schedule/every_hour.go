package schedule

import (
	"strconv"
	"time"

	"github.com/schedule-job/schedule-job-batch/external/tool"
)

type everyHour struct {
	Minute int `json:"minute"`
}

func everyHourToParams(params map[string]string) (*everyHour, error) {
	minute, minuteErr := strconv.Atoi(params["minute"])
	if minuteErr != nil {
		return nil, minuteErr
	}

	return &everyHour{
		Minute: minute,
	}, nil
}

func EveryHour(pivotTime time.Time, params map[string]string, _ interface{}) (*time.Time, error) {
	var payload, err = everyHourToParams(params)

	if err != nil {
		return nil, err
	}

	if pivotTime.Minute() > payload.Minute {
		pivotTime = pivotTime.Add(time.Hour)
	}

	newTime := tool.NewUTCDate(pivotTime.Year(), pivotTime.Month(), pivotTime.Day(), pivotTime.Hour(), payload.Minute)

	return &newTime, nil
}
