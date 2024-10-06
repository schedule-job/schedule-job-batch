package rule_based_replace

import (
	"errors"
	"strconv"
)

func ToTimestampAddMinute(params map[string]string, _options interface{}) (string, error) {
	if len(params) < 3 {
		return "", errors.New("hour or minute is empty")
	}

	timestamp, err := ToTimestamp(params, _options)
	if err != nil {
		return "", err
	}

	t, tErr := strconv.ParseInt(timestamp, 10, 64)
	min, minErr := strconv.ParseInt(params["add"], 10, 64)

	if tErr != nil || minErr != nil {
		return "", errors.New("hour or minute is empty")
	}

	return strconv.FormatInt(t+min*60*1000, 10), nil
}
