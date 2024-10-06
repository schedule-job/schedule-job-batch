package main

import (
	"fmt"
	"strings"

	parser "github.com/Sotaneum/go-args-parser"
	"github.com/gin-gonic/gin"
	ginsession "github.com/go-session/gin-session"
	"github.com/schedule-job/schedule-job-batch/internal/pg"
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

	if options.TrustedProxies != "" {
		trustedProxies := strings.Split(options.TrustedProxies, ",")
		router.SetTrustedProxies(trustedProxies)
	}

	router.POST("/api/v1/schedule/next", func(ctx *gin.Context) {
		// - 특정 아이템의 다음 실행일 API
		ctx.JSON(200, gin.H{"code": 200, "data": "ok"})
	})

	router.POST("/api/v1/request", func(ctx *gin.Context) {
		// - 1000개씩 읽어서 작업을 진행. 모든 요청이 끝나야 다음 1000개 작업
		ctx.JSON(200, gin.H{"code": 200, "data": "ok"})
	})

	router.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(404, gin.H{"code": 404, "message": "접근 할 수 없는 페이지입니다!"})
	})

	fmt.Println("Started Agent! on " + options.Port)

	router.Run(":" + options.Port)
}
