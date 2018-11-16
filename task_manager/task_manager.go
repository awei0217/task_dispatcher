package task_manager

import (
	"errors"
	"github.com/shotdog/quartz"
	"task_dispatcher/common"
	"task_dispatcher/domain"
	"task_dispatcher/execute"
	"time"
)

var (
	Error_quartz_null = errors.New("the quartz is null")
)

var tm *TaskManager
var QZ *quartz.Quartz

type TaskManager struct {
	qz *quartz.Quartz
}

func init() {
	if tm == nil {
		QZ = quartz.New()
		QZ.BootStrap()
		tm = &TaskManager{qz: QZ}
	}
}
func GetTaskManger() *TaskManager {
	return tm
}

func (tm *TaskManager) ActiveTask(tk *domain.Task) error {
	// 设置任务运行中
	j := &quartz.Job{Id: tk.Id, Name: tk.Name, Group: tk.GroupNo, Expression: tk.Cron, Params: tk.Param, Active: tk.IsActivity, JobFunc: tm.InvokeJob, Url: tk.Url}
	err := tm.qz.AddJob(j)
	if err == nil {
		tk.IsActivity = 1
		tk.UpdateTask()
	} else {
	}
	return err
}

func (tm *TaskManager) RemoveTask(taskId int) error {
	err := tm.qz.RemoveJob(taskId)
	return err
}

func (tm *TaskManager) StopTask(taskId int) error {
	err1 := tm.qz.RemoveJob(taskId)
	tk, _ := domain.FindTaskById(taskId)
	if tk == nil {
		return nil
	}
	tk.IsActivity = 2
	tk.UpdateTask()
	if err1 != nil {
		return err1
	}
	return nil
}
func (tm *TaskManager) UpdateTask(tk *domain.Task) error {
	if tk == nil {
		return nil
	}
	j := &quartz.Job{Id: tk.Id, Name: tk.Name, Group: tk.GroupNo, Expression: tk.Cron, Params: tk.Param, Active: tk.IsActivity, JobFunc: tm.InvokeJob, Url: tk.Url}
	err := tm.qz.ModifyJob(j)
	if err != nil {
		return err
	}
	return nil
}

func (this *TaskManager) InvokeJob(taskId int, url, params string, nextTime time.Time) {
	tk, err := domain.FindTaskById(taskId)
	if err != nil {
		common.GetLog().Errorln("根据ID查询任务异常 ", taskId, err)
		return
	}
	if tk.Status == 2 { // 停止状态
		GetTaskManger().StopTask(taskId)
		return
	}
	scheduler := &execute.Execute{}
	scheduler.ExecuteTask(tk)
}
