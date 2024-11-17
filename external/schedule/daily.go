package schedule

import (
	"strconv"
	"time"

	"github.com/schedule-job/schedule-job-batch/external/tool"
)

type daily struct {
	Hour   int `json:"hour"`
	Minute int `json:"minute"`
}

func dailyToParams(params map[string]string) (*daily, error) {
	hour, hourErr := strconv.Atoi(params["hour"])

	if hourErr != nil {
		return nil, hourErr
	}

	minute, minuteErr := strconv.Atoi(params["minute"])
	if minuteErr != nil {
		return nil, minuteErr
	}

	return &daily{
		Hour:   hour,
		Minute: minute,
	}, nil
}

func Daily(pivotTime time.Time, params map[string]string, _ interface{}) (*time.Time, error) {
	var payload, err = dailyToParams(params)

	if err != nil {
		return nil, err
	}

	if pivotTime.Hour()*60+pivotTime.Minute() > payload.Hour*60+payload.Minute {
		pivotTime = pivotTime.AddDate(0, 0, 1)
	}

	newTime := tool.NewUTCDate(pivotTime.Year(), pivotTime.Month(), pivotTime.Day(), payload.Hour, payload.Minute)

	return &newTime, nil
}
