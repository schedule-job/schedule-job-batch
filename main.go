package main

import (
	"fmt"
	"strings"
	"time"

	parser "github.com/Sotaneum/go-args-parser"
	"github.com/gin-gonic/gin"
	ginsession "github.com/go-session/gin-session"
	"github.com/schedule-job/schedule-job-batch/internal/core"
	"github.com/schedule-job/schedule-job-batch/internal/pg"
	"github.com/schedule-job/schedule-job-batch/internal/request"
	"github.com/schedule-job/schedule-job-batch/internal/rule_based_replace"
	"github.com/schedule-job/schedule-job-batch/internal/schedule"
	"github.com/schedule-job/schedule-job-batch/internal/tool"
)

type Options struct {
	Port           string
	PostgresSqlDsn string
	TrustedProxies string
	AgentUrl       string
}

var DEFAULT_OPTIONS = map[string]string{
	"PORT":             "8080",
	"POSTGRES_SQL_DSN": "",
	"TRUSTED_PROXIES":  "127.0.0.1",
	"AGENT_URL":        "",
}

func getOptions() *Options {
	rawOptions := parser.ArgsJoinEnv(DEFAULT_OPTIONS)

	options := new(Options)
	options.Port = rawOptions["PORT"]
	options.PostgresSqlDsn = rawOptions["POSTGRES_SQL_DSN"]
	options.TrustedProxies = rawOptions["TRUSTED_PROXIES"]
	options.AgentUrl = rawOptions["AGENT_URL"]

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
	if len(options.PostgresSqlDsn) == 0 {
		panic("not found 'AGENT_URL' options")
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
		payload := make(map[string]string)
		bindErr := ctx.BindJSON(&payload)

		if bindErr != nil {
			ctx.JSON(400, gin.H{"code": 400, "message": "payload 에러 : " + bindErr.Error()})
			return
		}

		result, err := core.GetNextSchedule(name, payload, time.Now().UTC())

		if err != nil {
			ctx.JSON(400, gin.H{"code": 400, "message": err.Error()})
			return
		}

		ctx.JSON(200, gin.H{"code": 200, "data": result})
	})

	router.POST("/api/v1/request/pre-next/:name", func(ctx *gin.Context) {
		name := ctx.Param("name")
		payload := make(map[string]interface{})
		bindErr := ctx.BindJSON(&payload)

		if bindErr != nil {
			ctx.JSON(400, gin.H{"code": 400, "message": "payload 에러 : " + bindErr.Error()})
			return
		}

		result, err := core.GetNextRequest(name, payload)

		if err != nil {
			ctx.JSON(400, gin.H{"code": 400, "message": err.Error()})
			return
		}

		ctx.JSON(200, gin.H{"code": 200, "data": request.NewRequestFormat(result)})
	})

	router.POST("/api/v1/schedule/next/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		result, err := core.GetNextScheduleByDatabase(id, database, time.Now().UTC())

		if err != nil {
			ctx.JSON(400, gin.H{"code": 400, "message": err.Error()})
			return
		}

		ctx.JSON(200, gin.H{"code": 200, "data": result})
	})

	router.POST("/api/v1/request/next/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		result, err := core.GetNextRequestByDatabase(id, database)

		if err != nil {
			ctx.JSON(400, gin.H{"code": 400, "message": err.Error()})
			return
		}

		ctx.JSON(200, gin.H{"code": 200, "data": request.NewRequestFormat(result)})
	})

	router.POST("/api/v1/progress/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		req, err := core.GetNextRequestByDatabase(id, database)

		if err != nil {
			ctx.JSON(400, gin.H{"code": 400, "message": err.Error()})
			return
		}

		items := []request.RequestInterface{}
		items = append(items, req)

		if core.ReqeustAgent(items, options.AgentUrl) != nil {
			ctx.JSON(400, gin.H{"code": 400, "message": err.Error()})
			return
		}

		ctx.JSON(200, gin.H{"code": 200, "data": "in progress"})
	})

	router.POST("/api/v1/progress", func(ctx *gin.Context) {
		ids, err := database.GetIdsByRequests()
		if err != nil {
			ctx.JSON(400, gin.H{"code": 400, "message": err.Error()})
			return
		}

		requests := []request.RequestInterface{}
		now := time.Now().UTC()
		pivotTime := tool.NewUTCDate(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute())
		errorCnt := 0
		for _, id := range ids {
			sch, schErr := core.GetNextScheduleByDatabase(id, database, pivotTime)
			if schErr != nil {
				errorCnt++
				continue
			}

			if pivotTime.Compare(*sch) == 0 {
				continue
			}

			req, reqErr := core.GetNextRequestByDatabase(id, database)

			if reqErr != nil {
				errorCnt++
				continue
			}

			requests = append(requests, req)
		}

		agentErr := core.ReqeustAgent(requests, options.AgentUrl)

		if agentErr != nil {
			ctx.JSON(400, gin.H{"code": 400, "message": agentErr.Error()})
			return
		}

		ctx.JSON(200, gin.H{"code": 200, "data": "in progress", "all": len(ids), "started": len(requests), "error": errorCnt})
	})

	router.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(404, gin.H{"code": 404, "message": "접근 할 수 없는 페이지입니다!"})
	})

	fmt.Println("Started Agent! on " + options.Port)

	router.Run(":" + options.Port)
}
