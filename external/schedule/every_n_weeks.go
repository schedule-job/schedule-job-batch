package schedule

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/schedule-job/schedule-job-batch/external/tool"
)

type everyNWeeks struct {
	StartDate     string `json:"startDate"`
	IntervalWeeks int    `json:"intervalWeeks"`
}

func everyNWeeksToParams(params map[string]string) (*everyNWeeks, error) {
	if params["startDate"] == "" {
		return nil, errors.New("StartDate is empty")
	}

	intervalWeeks, intervalWeeksErr := strconv.Atoi(params["intervalWeeks"])
	if intervalWeeksErr != nil {
		return nil, intervalWeeksErr
	}

	return &everyNWeeks{
		StartDate:     params["startDate"],
		IntervalWeeks: intervalWeeks,
	}, nil
}

func EveryNWeeks(pivotTime time.Time, params map[string]string, _ interface{}) (*time.Time, error) {
	var payload, err = everyNWeeksToParams(params)

	if err != nil {
		return nil, err
	}

	items := strings.Split(payload.StartDate, "-")

	if len(items) != 5 || payload.IntervalWeeks == 0 {
		return nil, errors.New("Invalid StartDate or IntervalWeeks")
	}

	numbers := []int{}

	for _, item := range items {
		number, err := strconv.Atoi(item)
		if err != nil {
			return nil, err
		}
		numbers = append(numbers, number)
	}

	start := tool.NewUTCDate(numbers[0], time.Month(numbers[1]), numbers[2], numbers[3], numbers[4])

	if pivotTime.Before(start) {
		return &start, nil
	}

	daysSinceStart := int(pivotTime.Sub(start).Hours()) / 24
	weeksSinceStart := daysSinceStart / 7
	if daysSinceStart%7 > 0 {
		weeksSinceStart += 1
	}
	nextCheckWeek := weeksSinceStart / payload.IntervalWeeks
	if weeksSinceStart%payload.IntervalWeeks > 0 {
		nextCheckWeek += 1
	}
	nextCheckWeek *= payload.IntervalWeeks
	nextCheckDate := start.AddDate(0, 0, nextCheckWeek*7)

	if nextCheckDate.After(pivotTime) || pivotTime.Equal(nextCheckDate) {
		return &nextCheckDate, nil
	}

	nextCheckDate = nextCheckDate.AddDate(0, 0, payload.IntervalWeeks*7)

	return &nextCheckDate, nil
}
