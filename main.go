package main

import (
	"fmt"
	"github.com/kataras/iris"
	"runtime"
)
import "./handler"
import "./utils"
import "./job"
import "./worker"
func main()  {
	runtime.GOMAXPROCS(runtime.NumCPU())
	app := iris.New()
	app.StaticWeb("/", utils.Config.String("webroot"))
	go job.InitJobManager()
	go worker.InitWorker()
	app.Post("/job/save",handler.SaveJob)
	app.Get("/job/delete",handler.DeleteJob)
	app.Get("/job/list",handler.JobList)
	app.Put("/job/kill",handler.Killjob)
	app.Get("/worker/list",handler.WorkerList)


	app.Run(iris.Addr(":"+utils.Config.String("HttpPort")))
	fmt.Println("服务器退出1")
}

