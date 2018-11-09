package controller

import (
	"github.com/kataras/iris/mvc"
	"strconv"
	"encoding/json"
	"github.com/kataras/iris/context"
	"task_dispatcher/domain"
	"task_dispatcher/common"
)

type RegisterCenterController struct {
	InterceptController
} 
/**
	跳转到任务列表页面
 */
func (registerCenterController *RegisterCenterController)GetRegisterCenterListPage()mvc.View{

	return mvc.View{
		Name:"/register_center/register-center-list",
	}
}

/**
	根据条件分页查询任务
 */
func (registerCenterController *RegisterCenterController)PostFindRegisterCenterList(ctx context.Context)mvc.Response{

	p,_ := strconv.Atoi(ctx.FormValue("page"))
	l,_ := strconv.Atoi(ctx.FormValue("limit"))
	isDeath,_ := strconv.Atoi(ctx.FormValue("isDeath"))

	rc := &domain.RegisterCenter{
		Ip:ctx.FormValue("ip"),
		IsDeath:isDeath,
	}

	page := &common.Page{
		Page:p,
		Limit:l,
	}
	page,err := domain.FindRegisterCenterByPage(page,rc)
	if(err != nil){
		page.Code = 500
		page.Message = "操作失败"

	}else{
		page.Code = 200
		page.Message = "操作成功"
	}
	resByte ,_ := json.Marshal(page)
	res := string(resByte)
	return mvc.Response{
		ContentType:"application/json;charset=UTF-8",
		Text:res,
	}
}