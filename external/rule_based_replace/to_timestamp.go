package rule_based_replace

import (
	"errors"
	"strconv"
	"time"

	"github.com/schedule-job/schedule-job-batch/external/tool"
)

type ToTimestampOption struct {
	Pivot *time.Time
}

func ToTimestamp(params map[string]string, _options interface{}) (string, error) {
	if len(params) < 2 {
		return "", errors.New("hour or minute is empty")
	}

	hour, hErr := strconv.Atoi(params["hour"])
	min, mErr := strconv.Atoi(params["minute"])

	if hErr != nil || mErr != nil {
		return "", errors.New("hour or minute is not a number")
	}

	var options ToTimestampOption
	if _options == nil {
		options = ToTimestampOption{}
	} else {
		var __options, ok = _options.(ToTimestampOption)
		if !ok {
			return "", errors.New("invalid options")
		}
		options = __options
	}

	var now time.Time

	if options.Pivot != nil {
		now = *options.Pivot
	} else {
		now = time.Now().UTC()
	}

	return strconv.FormatInt(tool.NewUTCDate(now.Year(), now.Month(), now.Day(), hour, min).UnixMilli(), 10), nil
}
