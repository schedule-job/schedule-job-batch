# Schedule Job Batch

[![Docker Image Build With Push](https://github.com/schedule-job/schedule-job-batch/actions/workflows/docker-image-build-push.yml/badge.svg)](https://github.com/schedule-job/schedule-job-batch/actions/workflows/docker-image-build-push.yml) [![Docker Pulls](https://img.shields.io/docker/pulls/sotaneum/schedule-job-batch?logoColor=fff&logo=docker)](https://hub.docker.com/r/sotaneum/schedule-job-batch) [![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/schedule-job/schedule-job-batch?logo=go&logoColor=fff)](https://go.dev/) [![Docker Image Size (tag)](https://img.shields.io/docker/image-size/sotaneum/schedule-job-batch/latest?logoColor=fff&logo=docker)](https://hub.docker.com/r/sotaneum/schedule-job-batch) [![postgresql](https://img.shields.io/badge/14_or_higher-blue?logo=postgresql&logoColor=fff&label=PostgreSQL)](https://www.postgresql.org/)

## API

### [POST] /api/v1/schedule/next/:name

- Request

  - Body

    ```json
    // /api/v1/schedule/next/everyHour
    {
      "minute": "32"
    }
    ```

- Response

  - Body

    ```json
    // /api/v1/schedule/next/everyHour
    {
      "code": 200,
      "data": "2024-10-06T08:32:00Z"
    }
    ```
