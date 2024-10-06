package main

import (
	"fmt"
	"strings"
	"time"

	parser "github.com/Sotaneum/go-args-parser"
	"github.com/gin-gonic/gin"
	ginsession "github.com/go-session/gin-session"
	"github.com/schedule-job/schedule-job-batch/internal/pg"
	"github.com/schedule-job/schedule-job-batch/internal/request"
	"github.com/schedule-job/schedule-job-batch/internal/rule_based_replace"
	"github.com/schedule-job/schedule-job-batch/internal/schedule"
)

type Options struct {
	Port           string
	PostgresSqlDsn string
	TrustedProxies string
}

var DEFAULT_OPTIONS = map[string]string{
	"PORT":             "8080",
	"POSTGRES_SQL_DSN": "",
	"TRUSTED_PROXIES":  "",
}

func getOptions() *Options {
	rawOptions := parser.ArgsJoinEnv(DEFAULT_OPTIONS)

	options := new(Options)
	options.Port = rawOptions["PORT"]
	options.PostgresSqlDsn = rawOptions["POSTGRES_SQL_DSN"]
	options.TrustedProxies = rawOptions["TRUSTED_PROXIES"]

	return options
}

func safeGo(f func()) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered from panic:", r)
			}
		}()
		f()
	}()
}

func main() {
	options := getOptions()
	if len(options.PostgresSqlDsn) == 0 {
		panic("not found 'POSTGRES_SQL_DSN' options")
	}
	if len(options.Port) == 0 {
		panic("not found 'PORT' options")
	}

	database := pg.New(options.PostgresSqlDsn)

	router := gin.Default()
	router.Use(ginsession.New())

	schedule.Scheduler.AddSchedule("daily", schedule.Daily, nil)
	schedule.Scheduler.AddSchedule("everyHour", schedule.EveryHour, nil)
	schedule.Scheduler.AddSchedule("everyMonth", schedule.EveryMonth, nil)
	schedule.Scheduler.AddSchedule("everyNWeeks", schedule.EveryNWeeks, nil)
	schedule.Scheduler.AddSchedule("everyWeeks", schedule.EveryWeeks, nil)
	schedule.Scheduler.AddSchedule("everyYear", schedule.EveryYear, nil)
	schedule.Scheduler.AddSchedule("weeksOfEveryMonth", schedule.WeeksOfEveryMonth, nil)

	rule_based_replace.Replacer.AddRule("toTimestamp", rule_based_replace.ToTimestamp, nil)
	rule_based_replace.Replacer.AddRule("toTimestampAddMinute", rule_based_replace.ToTimestampAddMinute, nil)

	request.Requester.AddRequest("defaultRequest", request.NewDefaultRequestByInterface, nil)

	if options.TrustedProxies != "" {
		trustedProxies := strings.Split(options.TrustedProxies, ",")
		router.SetTrustedProxies(trustedProxies)
	}

	router.POST("/api/v1/schedule/pre-next/:name", func(ctx *gin.Context) {
		name := ctx.Param("name")

		var check = schedule.Scheduler.CheckSupportedSchedule(name)
		if !check {
			ctx.JSON(400, gin.H{"code": 400, "message": "지원하지 않는 스케줄입니다."})
			return
		}

		payload := make(map[string]string)
		ctx.BindJSON(&payload)

		result, err := schedule.Scheduler.Schedule(time.Now().UTC(), name, payload)

		if err != nil {
			ctx.JSON(400, gin.H{"code": 400, "message": err.Error()})
			return
		}

		ctx.JSON(200, gin.H{"code": 200, "data": result})
	})

	router.POST("/api/v1/request/pre-next/:name", func(ctx *gin.Context) {
		name := ctx.Param("name")

		var check = request.Requester.CheckSupportedRequest(name)
		if !check {
			ctx.JSON(400, gin.H{"code": 400, "message": "지원하지 않는 요청입니다."})
			return
		}

		payload := make(map[string]interface{})
		ctx.BindJSON(&payload)

		result, err := request.Requester.Request(name, payload)

		if err != nil {
			ctx.JSON(400, gin.H{"code": 400, "message": err.Error()})
			return
		}

		ctx.JSON(200, gin.H{"code": 200, "data": request.NewRequestFormat(result)})
	})

	router.POST("/api/v1/schedule/next/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")

		data, err := database.GetSchedule(id)
		if err != nil {
			ctx.JSON(400, gin.H{"code": 400, "message": err.Error()})
			return
		}

		var check = schedule.Scheduler.CheckSupportedSchedule(data.Name)
		if !check {
			ctx.JSON(400, gin.H{"code": 400, "message": "지원하지 않는 스케줄입니다."})
			return
		}

		result, err := schedule.Scheduler.Schedule(time.Now().UTC(), data.Name, data.Payload)

		if err != nil {
			ctx.JSON(400, gin.H{"code": 400, "message": err.Error()})
			return
		}

		ctx.JSON(200, gin.H{"code": 200, "data": result})
	})

	router.POST("/api/v1/request/next/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")

		data, err := database.GetRequests(id)
		if err != nil {
			ctx.JSON(400, gin.H{"code": 400, "message": err.Error()})
			return
		}

		var check = request.Requester.CheckSupportedRequest(data.Name)
		if !check {
			ctx.JSON(400, gin.H{"code": 400, "message": "지원하지 않는 요청입니다."})
			return
		}

		result, err := request.Requester.Request(data.Name, data.Payload)

		if err != nil {
			ctx.JSON(400, gin.H{"code": 400, "message": err.Error()})
			return
		}

		ctx.JSON(200, gin.H{"code": 200, "data": request.NewRequestFormat(result)})
	})

	router.POST("/api/v1/progress/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		ctx.JSON(200, gin.H{"code": 200, "data": "ok"})
	})

	router.POST("/api/v1/progress", func(ctx *gin.Context) {

		ctx.JSON(200, gin.H{"code": 200, "data": "ok"})
	})

	router.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(404, gin.H{"code": 404, "message": "접근 할 수 없는 페이지입니다!"})
	})

	fmt.Println("Started Agent! on " + options.Port)

	router.Run(":" + options.Port)
}
