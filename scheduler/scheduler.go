package scheduler

import (
	"task_dispatcher/config"
	"task_dispatcher/domain"
	"task_dispatcher/task_manager"
	"github.com/goinggo/mapstructure"
	"time"
	"task_dispatcher/common"
	"strings"
)

func init() {

	rc,_:=domain.FindRegisterCenterByIp(common.GetIp())

	if rc.TaskSliceCollection == "" {
		return
	}
	newSls := getTaskSlice(rc.TaskSliceCollection)
	res, err := config.Conn.Table("task").
		Fields(domain.TASK_SELECT_FIELD).
		Where("status", 1).
		Where("task_type",1).
		Where("task_slice","in",newSls).
		Get()
	if err == nil {
		tasks := &[]domain.Task{}
		mapstructure.Decode(res, tasks)
		for _, tk := range *tasks {
			task_manager.GetTaskManger().ActiveTask(&tk)
		}
	}
}

func StartScheduler() {
	go ForFindExecuteTask()
}

func ForFindExecuteTask() {
	for {
		rc,_:=domain.FindRegisterCenterByIp(common.GetIp())
		if rc == nil {
			time.Sleep(10 * time.Second)
			continue
		}
		if rc.TaskSliceCollection == "" {
			time.Sleep(10 * time.Second)
			continue
		}
		newSls := getTaskSlice(rc.TaskSliceCollection)
		res, err := config.Conn.Table("task").
			Fields(domain.TASK_SELECT_FIELD).
			Where("status", 1). // 启动状态
			Where("task_type",1). // 定时任务
			Where("is_activity", 2). // 没有运行的
			Where("task_slice","in",newSls).
			Get()
		if err == nil {
			tasks := &[]domain.Task{}
			mapstructure.Decode(res, tasks)
			for _, tk := range *tasks {
				task_manager.GetTaskManger().ActiveTask(&tk)
			}
		}
		time.Sleep(10 * time.Second)
	}
}

func getTaskSlice(taskSlice string)[]interface{}  {
	deathTaskSlice := strings.Split(taskSlice, ",")
	newSls := make([]interface{}, len(deathTaskSlice))
	for i, v := range deathTaskSlice {
		newSls[i] = v
	}
	return newSls
}

