package controller

import (
	"github.com/kataras/iris"

)


// InterceptController is the user authentication controller, a custom shared controller.
type InterceptController struct {

}
// 在请求到来之前对请求进行拦截
// BeginRequest saves login state to the context, the user id.
func (i *InterceptController) BeginRequest(ctx iris.Context) {


}
// 在响应浏览器之前进行拦截
// EndRequest is here just to complete the BaseController
// in order to be tell iris to call the `BeginRequest` before the main method.
func (i *InterceptController) EndRequest(ctx iris.Context) {}



