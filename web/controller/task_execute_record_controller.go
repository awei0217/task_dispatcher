package controller

import (
	"encoding/json"
	"github.com/kataras/iris/context"
	"github.com/kataras/iris/mvc"
	"log"
	"strconv"
	"task_dispatcher/common"
	"task_dispatcher/domain"
)

type TaskExecuteRecordController struct {
	InterceptController
}

/**
跳转到任务列表页面
*/
func (taskExecuteRecordController *TaskExecuteRecordController) GetTaskExecuteRecordListPage() mvc.View {

	return mvc.View{
		Name: "/task_execute_record/task-execute-record-list",
	}
}

/**
根据条件分页查询任务
*/
func (taskExecuteRecordController *TaskExecuteRecordController) PostFindTaskExecuteRecordList(ctx context.Context) mvc.Response {
	status, _ := strconv.Atoi(ctx.FormValue("status"))
	p, _ := strconv.Atoi(ctx.FormValue("page"))
	l, _ := strconv.Atoi(ctx.FormValue("limit"))

	ter := &domain.TaskExecuteRecord{
		TaskNo:  ctx.FormValue("taskNo"),
		Status:  status,
		GroupNo: ctx.FormValue("groupNo"),
	}
	page := &common.Page{
		Page:  p,
		Limit: l,
	}
	page, err := domain.FindTaskExecuteRecordByPage(page, ter)
	if err != nil {
		log.Println(err)
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
