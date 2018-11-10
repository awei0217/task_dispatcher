package main

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"task_dispatcher/common"
	"task_dispatcher/net_server"
	"task_dispatcher/scheduler"
	"task_dispatcher/web/controller"
)

func main() {
	app := iris.New()
	// You got full debug messages, useful when using MVC and you want to make
	// sure that your code is aligned with the Iris' MVC Architecture.
	app.Logger().SetLevel("error")

	// Load the template files.
	tmplate := iris.HTML("./web/views", ".html").Layout("shared/layout.html").Reload(true)

	app.RegisterView(tmplate)

	app.StaticWeb("/public", "./web/public")

	app.OnAnyErrorCode(func(ctx iris.Context) {
		ctx.ViewData("Message", ctx.Values().
			GetStringDefault("message", "The page you're looking for doesn't exist"))
		ctx.View("shared/error.html")
	})
	scheduler.StartScheduler()
	net_server.InitServer()
	mvc.New(app.Party("/")).Handle(new(controller.MainController))
	mvc.New(app.Party("/register/center")).Handle(new(controller.RegisterCenterController))
	mvc.New(app.Party("/task")).Handle(new(controller.TaskController))
	mvc.New(app.Party("/execute/record")).Handle(new(controller.TaskExecuteRecordController))
	mvc.New(app.Party("/monitor")).Handle(new(controller.MonitorController))
	iris.RegisterOnInterrupt(func() {
		common.GetLog().Println("程序异常终止")
		net_server.StopServer()
	})
	app.Run(
		// Starts the web server at localhost:8080
		iris.Addr("localhost:8080"),
		// Disables the updater.
		iris.WithoutVersionChecker,
		// Ignores err server closed log when CTRL/CMD+C pressed.
		iris.WithoutServerError(iris.ErrServerClosed),
		// Enables faster json serialization and more.
		iris.WithOptimizations,
		iris.WithConfiguration(iris.Configuration{
			DisableInterruptHandler: true,
		}),
	)
}



