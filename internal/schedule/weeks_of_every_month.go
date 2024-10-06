package schedule

import (
	"errors"
	"strconv"
	"time"

	"github.com/schedule-job/schedule-job-batch/internal/tool"
)

type weeksOfEveryMonth struct {
	MonthlyWeekNumbers []int `json:"monthlyWeekNumbers"`
	Weekdays           []int `json:"weekdays"`
	Hour               int   `json:"hour"`
	Minute             int   `json:"minute"`
}

type weeksOfEveryMonthOption struct {
	isNotCheckNextMonth bool
}

func weeksOfEveryMonthToParams(params map[string]string) (*weeksOfEveryMonth, error) {
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

	if params["monthlyWeekNumbers"] == "" {
		return nil, errors.New("monthlyWeekNumbers is empty")
	}

	monthlyWeekNumbers := tool.ConvertToInArray(params["monthlyWeekNumbers"])

	return &weeksOfEveryMonth{
		MonthlyWeekNumbers: monthlyWeekNumbers,
		Weekdays:           weekdays,
		Hour:               hour,
		Minute:             minute,
	}, nil
}

func WeeksOfEveryMonth(pivotTime time.Time, params map[string]string, _options interface{}) (*time.Time, error) {
	var payload, err = weeksOfEveryMonthToParams(params)

	if err != nil {
		return nil, err
	}

	if len(payload.MonthlyWeekNumbers) == 0 || len(payload.Weekdays) == 0 {
		return nil, errors.New("Invalid MonthlyWeekNumbers or Weekdays")
	}

	startOfMonth := tool.NewUTCDate(pivotTime.Year(), pivotTime.Month(), 1, payload.Hour, payload.Minute)
	startDayOfWeek := int(startOfMonth.Weekday())

	for _, week := range payload.MonthlyWeekNumbers {
		for _, dayOfWeek := range payload.Weekdays {
			if week == 1 && dayOfWeek < startDayOfWeek {
				continue
			}

			dayDiff := (week-1)*7 + dayOfWeek - startDayOfWeek

			nextDate := startOfMonth.AddDate(0, 0, dayDiff)
			if (nextDate.After(pivotTime) && nextDate.Month() == pivotTime.Month()) || pivotTime.Equal(nextDate) {
				return &nextDate, nil
			}
		}
	}

	var options weeksOfEveryMonthOption

	if _options == nil {
		options = weeksOfEveryMonthOption{}
	} else {
		__options, check := _options.(weeksOfEveryMonthOption)

		if !check {
			return nil, errors.New("Invalid options")
		}

		options = __options
	}

	if options.isNotCheckNextMonth {
		return nil, errors.New("no valid date")
	}

	options.isNotCheckNextMonth = true

	return WeeksOfEveryMonth(startOfMonth.AddDate(0, 1, 0), params, options)
}
