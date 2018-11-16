package controller

import (
	"encoding/json"
	"github.com/kataras/iris/context"
	"github.com/kataras/iris/mvc"
	"log"
	"strconv"
	"task_dispatcher/common"
	"task_dispatcher/domain"
	"time"
)

type MonitorController struct {
	InterceptController
}

/**
跳转到性能监控页面
*/
func (monitorController *MonitorController) GetPerformanceMonitorPage() mvc.View {
	rc, _ := domain.FindAll()
	result := make(map[string]string)
	for i := 0; i < len(*rc); i++ {
		result[(*rc)[i].Ip] = (*rc)[i].Ip
	}
	return mvc.View{
		Name: "/monitor/performance-monitor",
		Data: result,
	}
}

/**
跳转到任务监控页面
*/
func (monitorController *MonitorController) GetTaskMonitorPage() mvc.View {
	rc, _ := domain.FindAll()
	result := make(map[string]string)
	for i := 0; i < len(*rc); i++ {
		result[(*rc)[i].Ip] = (*rc)[i].Ip
	}
	return mvc.View{
		Name: "/monitor/task-monitor",
		Data: result,
	}
}

/**
获取性能监控信息
*/
func (monitorController *MonitorController) PostFindPerformanceMonitorInfo(ctx context.Context) mvc.Response {
	ip := ctx.FormValue("ipAddress")
	timeBeginAndEnd, _ := strconv.Atoi(ctx.FormValue("timeBeginAndEnd"))
	limit := 0
	startTime := ""
	endTime := ""
	if timeBeginAndEnd == 1 {
		limit = 360
	} else if timeBeginAndEnd == 5 {
		limit = 1800
	} else if timeBeginAndEnd == 8 {
		limit = 3000
	} else if timeBeginAndEnd == 0 { // 今天
		timeDay := time.Now().Format(common.YYYY_MM_DD)
		startTime = timeDay + " 00:00:00"
		endTime = timeDay + " 23:59:59"
	} else if timeBeginAndEnd == -1 { // 昨天
		timeDay := time.Now().AddDate(0, 0, -1).Format(common.YYYY_MM_DD)
		startTime = timeDay + " 00:00:00"
		endTime = timeDay + " 23:59:59"
	}
	result := make(map[string]interface{})
	if startTime != "" {
		pms, err := domain.FindPerformanceMonitorByIpAndTime(ip, startTime, endTime)
		if err != nil {
			log.Println("查询监控信息异常")
		}
		length := len(*pms)
		updateTimes := make([]string, length, length)
		cpus := make([]string, length, length)
		memons := make([]string, length, length)
		for i, j := len(*pms), 0; i > 0; i-- {
			updateTimes[j] = (*pms)[i-1].UpdateTime
			cpus[j] = (*pms)[i-1].UseCpu
			memons[j] = (*pms)[i-1].UseMemory
			j++
		}
		result["updateTimes"] = updateTimes
		result["cpus"] = cpus
		result["memorys"] = memons
	} else {
		pms, err := domain.FindPerformanceMonitor(ip, limit)
		if err != nil {
			log.Println("查询监控信息异常")
		}
		length := len(*pms)

		updateTimes := make([]string, length, length)
		cpus := make([]string, length, length)
		memons := make([]string, length, length)
		for i, j := len(*pms), 0; i > 0; i-- {
			updateTimes[j] = (*pms)[i-1].UpdateTime
			cpus[j] = (*pms)[i-1].UseCpu
			memons[j] = (*pms)[i-1].UseMemory
			j++
		}
		result["updateTimes"] = updateTimes
		result["cpus"] = cpus
		result["memorys"] = memons
	}
	res, _ := json.Marshal(result)
	return mvc.Response{
		ContentType: "application/json;charset=UTF-8",
		Text:        string(res),
	}
}

func (monitorController *MonitorController) PostFindTaskMonitorInfo(ctx context.Context) mvc.Response {
	// 总任务
	count, _ := domain.CountTask()
	// 执行中的任务
	countExecuting, _ := domain.CountExecutingTask()
	// 定时任务
	countCron, _ := domain.CountCronTask()
	// 普通任务
	countOne, _ := domain.CountOneTask()
	countSlice := make([]int64, 4, 4)
	countSlice[0] = count
	countSlice[1] = countExecuting
	countSlice[2] = countCron
	countSlice[3] = countOne
	result := make(map[string]interface{})
	result["count"] = countSlice
	res, _ := json.Marshal(result)
	return mvc.Response{
		ContentType: "application/json;charset=UTF-8",
		Text:        string(res),
	}

}
