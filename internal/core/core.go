package core

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/schedule-job/schedule-job-batch/internal/request"
	"github.com/schedule-job/schedule-job-batch/internal/schedule"
	"github.com/schedule-job/schedule-job-database/pg"
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
		return nil, errors.New("요청 정보가 없습니다.")
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
	data, err := database.SelectAction(id)

	if err != nil {
		return nil, errors.New("요청 정보가 없습니다.")
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
		return errors.New("API 요청 포맷에 문제가 발생했습니다. 지속된 이슈가 발생하는 경우 관리자에게 문의하세요.")
	}
	api.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	result, clientErr := client.Do(api)

	if clientErr != nil {
		return errors.New("Agent 서버와 연결을 실패했습니다.")
	}

	defer result.Body.Close()

	if result.StatusCode != 200 {
		return errors.New("Agent 서버 요청에 실패했습니다.")
	}

	_, readErr := io.ReadAll(result.Body)

	if readErr != nil {
		return readErr
	}

	return nil
}
