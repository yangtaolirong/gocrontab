package handler

import (
	"../job"
	"../utils"
	//"../worker"
	"fmt"
	"github.com/kataras/iris/context"
	"net/http"
)


func SaveJob(context context.Context)  {
	//path := context.Path()
	//2.Json数据解析
	var j job.Job
	//context.ReadJSON()
	if err := context.ReadJSON(&j); err != nil {
		panic(err.Error())
	}
	fmt.Println()
	///fmt.Println("j=================",j)
	//fmt.Println(job.JM)
	oldjob,err:=job.JM.JobSave(&j)
	var resdata []byte
	if err!=nil{
		fmt.Println(err)
		resdata=utils.BuildRes(http.StatusInternalServerError,err.Error(),nil)
	}
   resdata=utils.BuildRes(http.StatusOK,"ok",oldjob)
	//输出：Received: main.Person{Name:"davie", Age:28}
	context.Write(resdata)
}
func DeleteJob(context context.Context)  {
	name:= context.URLParam("name") //
	fmt.Println(name,"=====")
	key:=utils.JOB_SAVE_DIR+name
	oldjob,err:=job.JM.JobDelete(key)
	var resdata []byte
	if err!=nil{
		fmt.Println(err)
		resdata=utils.BuildRes(http.StatusInternalServerError,err.Error(),nil)
	}
	resdata=utils.BuildRes(http.StatusOK,"ok",oldjob)

	//输出：Received: main.Person{Name:"davie", Age:28}
	context.Write(resdata)
}
func JobList(context context.Context)  {
	jobs,err:=job.JM.JobList()
	var resdata []byte
	if err!=nil{
		fmt.Println(err)
		resdata=utils.BuildRes(http.StatusInternalServerError,err.Error(),nil)
	}
	resdata=utils.BuildRes(http.StatusOK,"ok",jobs)

	//输出：Received: main.Person{Name:"davie", Age:28}
	context.Write(resdata)
}
func Killjob(context context.Context)  {
	name := context.URLParam("name")
	err:=job.JM.JobKill(name)
	var resdata []byte
	if err!=nil{
		fmt.Println(err)
		resdata=utils.BuildRes(http.StatusInternalServerError,err.Error(),nil)
	}
	resdata=utils.BuildRes(http.StatusOK,"ok",nil)

	//输出：Received: main.Person{Name:"davie", Age:28}
	context.Write(resdata)
}
func WorkerList(context context.Context)  {
	ips,err:=job.G_workerMgr.ListWorkers()
	var resdata []byte
	if err!=nil{
		fmt.Println(err)
		resdata=utils.BuildRes(http.StatusInternalServerError,err.Error(),nil)
	}
	resdata=utils.BuildRes(http.StatusOK,"ok",ips)
	context.Write(resdata)

}
