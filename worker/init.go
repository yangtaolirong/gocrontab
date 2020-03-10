package worker

import "fmt"
import "../job"
func InitWorker()  {
	job.InitWorkerMgr()
	err:=InitExecutor()
	if err!=nil{
		fmt.Println(err)
	}
	err=InitRegister()
	if err!=nil{
		fmt.Println(err)
	}
	err=InitScheduler()
	if err!=nil{
		fmt.Println(err)
	}
	err=InitJobMgr()
	if err!=nil{
		fmt.Println(err)
	}
}
