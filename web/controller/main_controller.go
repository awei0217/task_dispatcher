package controller

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
)

type MainController struct {
	InterceptController
	Ctx iris.Context
}

var indexView = mvc.View{
	Name: "index",
}

func (mainController *MainController) Get() mvc.View {

	return indexView

}
