# Schedule Job Batch

[![Docker Image Build With Push](https://github.com/schedule-job/schedule-job-batch/actions/workflows/docker-image-build-push.yml/badge.svg)](https://github.com/schedule-job/schedule-job-batch/actions/workflows/docker-image-build-push.yml) [![Docker Pulls](https://img.shields.io/docker/pulls/sotaneum/schedule-job-batch?logoColor=fff&logo=docker)](https://hub.docker.com/r/sotaneum/schedule-job-batch) [![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/schedule-job/schedule-job-batch?logo=go&logoColor=fff)](https://go.dev/) [![Docker Image Size (tag)](https://img.shields.io/docker/image-size/sotaneum/schedule-job-batch/latest?logoColor=fff&logo=docker)](https://hub.docker.com/r/sotaneum/schedule-job-batch) [![postgresql](https://img.shields.io/badge/14_or_higher-blue?logo=postgresql&logoColor=fff&label=PostgreSQL)](https://www.postgresql.org/)

- Data is processed based on DB and delivered to [Agent](https://github.com/schedule-job/schedule-job-agent).
- DB기반으로 데이터를 가공하여 [Agent](https://github.com/schedule-job/schedule-job-agent)에 전달합니다.

## Addon

### Schedule
- Calculate the date on which it will next run.
- 다음 일정을 계산합니다.


### Request
- Provides the ability to request API.
- API 요청할 수 있는 기능을 제공합니다.

#### Rule Based Replace
- Supports Body function.
- Body함수를 지원합니다.
- `to_timestamp`
  - Returns a timestamp based on when the task runs.
  - 작업이 실행되는 시점을 기준으로 timestamp를 반환합니다.
  - input
    ```text
    [:toTimestamp(hour=4, minute=4):]
    ```
  - output
    ```text
    1609473840000
    ```
- `to_timestamp_add_minute`
  - Returns a timestamp by adding minutes to when the task runs.
  - 작업이 실행되는 시점에 분을 추가하여 timestamp를 반환합니다.
  - input
    ```text
    [:toTimestampAddMinute(hour=4, minute=4, add=4):]
    ```
  - output
    ```text
    1609474080000
    ```

### Progress
- Requests are performed based on the DB based on the time of request.
- 요청 시점을 기준으로 DB 기준으로 Reqeust를 수행합니다.

## API

### [POST] /api/v1/schedule/pre-next/:name

- Example

  - `/api/v1/schedule/pre-next/everyHour`
  - Request

    - Body

      ```json
      {
        "minute": "32"
      }
      ```

  - Response

    - Body

      ```json
      {
        "code": 200,
        "data": "2024-10-06T08:32:00Z"
      }
      ```

### [POST] /api/v1/request/pre-next/:name

- Example

  - `/api/v1/request/pre-next/defaultRequest`
  - Request

    - Body

      ```json
      {
        "id": "1234",
        "url": "https://localhost:8080",
        "method": "GET",
        "body": "[:toTimestamp(hour=10, minute=10):]",
        "headers": {}
      }
      ```

  - Response

    - Body

      ```json
      {
        "code": 200,
        "data": {
          "id": "1234",
          "url": "https://localhost:8080",
          "method": "GET",
          "body": "1728209400000",
          "headers": {}
        }
      }
      ```

### [POST] /api/v1/schedule/next/:id

- Example

  - `/api/v1/schedule/next/12345678-1234-5678-1234-567812345678`
  - Response

    - Body

      ```json
      {
        "code": 200,
        "data": "2024-10-06T08:32:00Z"
      }
      ```

### [POST] /api/v1/request/next/:id

- Example

  - `/api/v1/request/next/12345678-1234-5678-1234-567812345678`
  - Response

    - Body

      ```json
      {
        "code": 200,
        "data": {
          "id": "12345678-1234-5678-1234-567812345678",
          "url": "https://localhost:8080",
          "method": "GET",
          "body": "1728209400000",
          "headers": {}
        }
      }
      ```

### [POST] /api/v1/progress/:id

- Example

  - `/api/v1/progress/12345678-1234-5678-1234-567812345678`
  - Response

    - Body

      ```json
      {
        "code": 200,
        "data": "in progress"
      }
      ```

### [POST] /api/v1/progress

- Example

  - `/api/v1/progress`
  - Response

    - Body

      ```json
      {
        "code": 200
        "data": "in progress"
      }
      ```
