package core

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/schedule-job/schedule-job-batch/internal/pg"
	"github.com/schedule-job/schedule-job-batch/internal/request"
	"github.com/schedule-job/schedule-job-batch/internal/schedule"
)

func GetNextSchedule(name string, payload map[string]string, pivotTime time.Time) (*time.Time, error) {
	var check = schedule.Scheduler.CheckSupportedSchedule(name)
	if !check {
		return nil, errors.New("지원하지 않는 스케줄입니다.")
	}

	return schedule.Scheduler.Schedule(pivotTime, name, payload)
}

func GetNextScheduleByDatabase(id string, database *pg.PostgresSQL, pivotTime time.Time) (*time.Time, error) {
	data, err := database.GetSchedule(id)

	if err != nil {
		return nil, err
	}

	return GetNextSchedule(data.Name, data.Payload, pivotTime)
}

func GetNextRequest(name string, payload map[string]interface{}) (request.RequestInterface, error) {
	var check = request.Requester.CheckSupportedRequest(name)
	if !check {
		return nil, errors.New("지원하지 않는 요청입니다.")
	}

	return request.Requester.Request(name, payload)
}

func GetNextRequestByDatabase(id string, database *pg.PostgresSQL) (request.RequestInterface, error) {
	data, err := database.GetRequest(id)

	if err != nil {
		return nil, err
	}

	return GetNextRequest(data.Name, data.Payload)
}

func ReqeustAgent(reqs []request.RequestInterface, agentUrl string) error {
	formats := []request.RequestFormat{}

	for _, item := range reqs {
		format := request.NewRequestFormat(item)
		formats = append(formats, format)
	}

	jsonData, err := json.Marshal(formats)
	if err != nil {
		return err
	}

	api, createAPIErr := http.NewRequest("POST", agentUrl+"/api/v1/request", bytes.NewBuffer(jsonData))
	if createAPIErr != nil {
		return createAPIErr
	}
	api.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	result, clientErr := client.Do(api)

	if clientErr != nil {
		return clientErr
	}

	defer result.Body.Close()

	_, readErr := io.ReadAll(result.Body)

	if readErr != nil {
		return readErr
	}

	return nil
}
