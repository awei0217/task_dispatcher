package controller

import (
	"task_dispatcher/common"
	"task_dispatcher/execute"
	"task_dispatcher/domain"
	"encoding/json"
	"github.com/kataras/iris/context"
	"github.com/kataras/iris/mvc"
	"strconv"
	"task_dispatcher/task_manager"
	"strings"
)

type TaskController struct {
	InterceptController
}




/**
跳转到任务列表页面
*/
func (taskController *TaskController) GetTaskListPage() mvc.View {

	return mvc.View{
		Name: "/task/task-list",
	}
}

/**
跳转到任务添加页面
*/
func (taskController *TaskController) GetAddTaskPage() mvc.View {

	return mvc.View{
		Name: "/task/add-task",
	}
}

/**
跳转到任务修改页面
*/
func (taskController *TaskController) GetUpdateTaskPage(ctx context.Context) mvc.View {
	taskId, _ := strconv.Atoi(ctx.FormValue("id"))
	tk, _ := domain.FindTaskById(taskId)
	return mvc.View{
		Name: "/task/update-task",
		Data: tk,
	}
}

/**
根据条件分页查询任务
*/
func (taskController *TaskController) PostFindTaskList(ctx context.Context) mvc.Response {

	tk := &domain.Task{
		TaskNo:  ctx.FormValue("taskNo"),
		Name:    ctx.FormValue("name"),
		GroupNo: ctx.FormValue("groupNo"),
	}
	p, _ := strconv.Atoi(ctx.FormValue("page"))
	l, _ := strconv.Atoi(ctx.FormValue("limit"))
	page := &common.Page{
		Page:  p,
		Limit: l,
	}
	page, err := domain.FindTaskByPage(page, tk)
	if err != nil {
		page.Code = 500
		page.Message = "操作失败"

	} else {
		page.Code = 200
		page.Message = "操作成功"
	}
	resByte, _ := json.Marshal(page)
	res := string(resByte)
	return mvc.Response{
		ContentType: "application/json;charset=UTF-8",
		Text:        res,
	}
}

/**
添加任务
*/
func (taskController *TaskController) PostAddTask(ctx context.Context) mvc.Response {
	isActivity := 2 // 默认不是运行中
	taskType, _ := strconv.Atoi(ctx.FormValue("taskType"))
	if taskType == 2{ // 手动任务
		isActivity = 3
	}
	status, _ := strconv.Atoi(ctx.FormValue("status"))
	isRegion, _ := strconv.Atoi(ctx.FormValue("isRegion"))
	isRecordLog, _ := strconv.Atoi(ctx.FormValue("isRecordLog"))
	concurrencyNum, _ := strconv.Atoi(ctx.FormValue("concurrencyNum"))
	requestMethod, _ := strconv.Atoi(ctx.FormValue("requestMethod"))

	tk := &domain.Task{
		Name:           ctx.FormValue("name"),
		GroupNo:        ctx.FormValue("groupNo"),
		Mail:           ctx.FormValue("mail"),
		Url:            ctx.FormValue("url"),
		Param:          ctx.FormValue("param"),
		Cron:           ctx.FormValue("cron"),
		Description:    ctx.FormValue("description"),
		TaskType:       taskType,
		Status:         status,
		IsActivity:     isActivity,
		IsRecordLog:    isRecordLog,
		ConcurrencyNum: concurrencyNum,
		CreateUser:     "system",
		UpdateUser:     "system",
		IsRegion:		isRegion,
		RequestMethod:  requestMethod,
	}
	_, err := tk.AddTask()
	if err != nil {
		result, _ := json.Marshal(&common.ControllerResponse{Code: 500, Message: err.Error()})
		return mvc.Response{
			ContentType: "application/json;charset=UTF-8",
			Text:        string(result),
		}
	} else {
		result, _ := json.Marshal(&common.ControllerResponse{Code: 200, Message: "操作成功"})
		return mvc.Response{
			ContentType: "application/json;charset=UTF-8",
			Text:        string(result),
		}
	}
}

/**
修改任务
*/
func (taskController *TaskController) PostUpdateTask(ctx context.Context) mvc.Response {

	taskType, _ := strconv.Atoi(ctx.FormValue("taskType"))
	isRegion, _ := strconv.Atoi(ctx.FormValue("isRegion"))
	status, _ := strconv.Atoi(ctx.FormValue("status"))
	isRecordLog, _ := strconv.Atoi(ctx.FormValue("isRecordLog"))
	concurrencyNum, _ := strconv.Atoi(ctx.FormValue("concurrencyNum"))
	requestMethod, _ := strconv.Atoi(ctx.FormValue("requestMethod"))
	taskId, _ := strconv.Atoi(ctx.FormValue("id"))
	cron := ctx.FormValue("cron")
	tk, _ := domain.FindTaskById(taskId)
	oldCron := tk.Cron
	oldTaskType := tk.TaskType
	tk.Name = ctx.FormValue("name")
	tk.GroupNo = ctx.FormValue("groupNo")
	tk.Mail = ctx.FormValue("mail")
	tk.Url = ctx.FormValue("url")
	tk.Param = ctx.FormValue("param")
	tk.Cron = cron
	tk.Description = ctx.FormValue("description")
	if oldTaskType != taskType && taskType == 1{ // 说明将手动任务改成自动任务
		tk.IsActivity = 2
	}
	if taskType == 2 && oldTaskType == 1{ // 说明将自动任务改成手动任务
		tk.IsActivity = 3
	}
	tk.TaskType = taskType
	tk.Status = status
	tk.IsRecordLog = isRecordLog
	tk.ConcurrencyNum = concurrencyNum
	tk.IsRegion = isRegion
	tk.RequestMethod =requestMethod
	_, err := tk.UpdateTask()

	// 如果任务是定时任务  并且 是运行状态的 并且 是改了cron表达式的
	if tk.TaskType == 1 && tk.Status == 1 && tk.IsActivity == 1 && oldCron != cron {
		task_manager.GetTaskManger().UpdateTask(tk)
	}
	if err != nil {
		result, _ := json.Marshal(&common.ControllerResponse{Code: 500, Message: err.Error()})
		return mvc.Response{
			ContentType: "application/json;charset=UTF-8",
			Text:        string(result),
		}
	} else {
		result, _ := json.Marshal(&common.ControllerResponse{Code: 200, Message: "操作成功"})
		return mvc.Response{
			ContentType: "application/json;charset=UTF-8",
			Text:        string(result),
		}
	}
}

/**
删除任务
*/
func (taskController *TaskController) PostDeleteTask(ctx context.Context) mvc.Response {

	taskId, _ := strconv.Atoi(ctx.FormValue("id"))

	_, err := domain.DeleteTask(taskId)
	if err != nil {
		result, _ := json.Marshal(&common.ControllerResponse{Code: 500, Message: err.Error()})
		return mvc.Response{
			ContentType: "application/json;charset=UTF-8",
			Text:        string(result),
		}
	} else {
		result, _ := json.Marshal(&common.ControllerResponse{Code: 200, Message: "操作成功"})
		return mvc.Response{
			ContentType: "application/json;charset=UTF-8",
			Text:        string(result),
		}
	}
}

/**
根据ID执行定时任务
*/
func (taskController *TaskController) PostExecuteTimingTask(ctx context.Context) mvc.Response {
	ids := ctx.FormValue("ids")
	idArray := strings.Split(ids,",")
	for i:=0;i<len(idArray);i++ {
		id,_  := strconv.Atoi(idArray[i])
		tk, err := domain.FindTaskById(id)
		if err != nil {
			continue
		}
		// 定时任务
		if tk.TaskType == 1 {
			tk.Status = 1
			tk.UpdateTask()

			// 普通一次性任务
		} else {
			execute := &execute.Execute{}
			execute.ExecuteTask(tk)

		}
	}
	return mvc.Response{
		ContentType: "application/json;charset=UTF-8",
		Text:        "启动成功",
	}
}

/**
根据ID停止定时任务
*/
func (taskController *TaskController) PostExecuteStopTask(ctx context.Context) mvc.Response {
	ids := ctx.FormValue("ids")
	idArray := strings.Split(ids,",")
	for i:=0;i<len(idArray);i++ {
		id,_  := strconv.Atoi(idArray[i])
		tk, err := domain.FindTaskById(id)
		if err != nil {
			continue
		}
		// 定时任务
		if tk.TaskType == 1 {
			tk.Status = 2
			tk.UpdateTask()
		}
	}
	return mvc.Response{
		ContentType: "application/json;charset=UTF-8",
		Text:        "任务停止成功",
	}
}
