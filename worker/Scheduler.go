package worker

import (
	"../job"
	"os"
	"time"
	"fmt"
	"../utils"
)

// 任务调度
type Scheduler struct {
	jobEventChan chan *job.JobEvent	//  etcd任务事件队列
	jobPlanTable map[string]*job.JobSchedulePlan // 任务调度计划表
	jobExecutingTable map[string]*job.JobExecuteInfo // 任务执行表
	jobResultChan chan *job.JobExecuteResult	// 任务结果队列
}

var (
	G_scheduler *Scheduler
)

// 处理任务事件
func (scheduler *Scheduler) handleJobEvent(jobEvent *job.JobEvent) {
	var (
		jobSchedulePlan *job.JobSchedulePlan
		jobExecuteInfo *job.JobExecuteInfo
		jobExecuting bool
		jobExisted bool
		err error
	)
	switch jobEvent.EventType {
	case utils.JOB_EVENT_SAVE:	// 保存任务事件
		if jobSchedulePlan, err = job.BuildJobSchedulePlan(jobEvent.Job); err != nil {
			return
		}
		scheduler.jobPlanTable[jobEvent.Job.Name] = jobSchedulePlan
	case  utils.JOB_EVENT_DELETE: // 删除任务事件
		if jobSchedulePlan, jobExisted = scheduler.jobPlanTable[jobEvent.Job.Name]; jobExisted {
			delete(scheduler.jobPlanTable, jobEvent.Job.Name)
		}
	case  utils.JOB_EVENT_KILL: // 强杀任务事件
		// 取消掉Command执行, 判断任务是否在执行中
		if jobExecuteInfo, jobExecuting = scheduler.jobExecutingTable[jobEvent.Job.Name]; jobExecuting {
			jobExecuteInfo.CancelFunc()	// 触发command杀死shell子进程, 任务得到退出
		}
	}
}

// 尝试执行任务
func (scheduler *Scheduler) TryStartJob(jobPlan *job.JobSchedulePlan) {
	// 调度 和 执行 是2件事情
	var (
		jobExecuteInfo *job.JobExecuteInfo
		jobExecuting bool
	)

	// 执行的任务可能运行很久, 1分钟会调度60次，但是只能执行1次, 防止并发！

	// 如果任务正在执行，跳过本次调度
	if jobExecuteInfo, jobExecuting = scheduler.jobExecutingTable[jobPlan.Job.Name]; jobExecuting {
		// fmt.Println("尚未退出,跳过执行:", jobPlan.Job.Name)
		return
	}

	// 构建执行状态信息
	jobExecuteInfo = job.BuildJobExecuteInfo(jobPlan)

	// 保存执行状态
	scheduler.jobExecutingTable[jobPlan.Job.Name] = jobExecuteInfo
	// 执行任务
	fmt.Println("执行任务:", jobExecuteInfo.Job.Name, jobExecuteInfo.PlanTime, jobExecuteInfo.RealTime)
	G_executor.ExecuteJob(jobExecuteInfo)
}

// 重新计算任务调度状态
func (scheduler *Scheduler) TrySchedule() (scheduleAfter time.Duration) {
	var (
		jobPlan *job.JobSchedulePlan
		now time.Time
		nearTime *time.Time
	)

	// 如果任务表为空话，随便睡眠多久
	if len(scheduler.jobPlanTable) == 0 {
		scheduleAfter = 1 * time.Second
		return
	}

	// 当前时间
	now = time.Now()

	// 遍历所有任务
	for _, jobPlan = range scheduler.jobPlanTable {
		if jobPlan.NextTime.Before(now) || jobPlan.NextTime.Equal(now) {
			scheduler.TryStartJob(jobPlan)
			jobPlan.NextTime = jobPlan.Expr.Next(now) // 更新下次执行时间
		}

		// 统计最近一个要过期的任务时间
		if nearTime == nil || jobPlan.NextTime.Before(*nearTime) {
			nearTime = &jobPlan.NextTime
		}
	}
	// 下次调度间隔（最近要执行的任务调度时间 - 当前时间）
	scheduleAfter = (*nearTime).Sub(now)
	return
}

// 处理任务结果
func (scheduler *Scheduler) handleJobResult(result *job.JobExecuteResult) {
	var (
		//jobLog *common.JobLog
	)
	// 删除执行状态
	delete(scheduler.jobExecutingTable, result.ExecuteInfo.Job.Name)
	file1,err:=os.OpenFile("logs/crontab.logs",os.O_APPEND|os.O_RDWR,0664)
	if err!=nil{
		file1,err:= os.Create("logs/crontab.logs")
		if err!=nil{
			fmt.Println(err)
		}
		defer file1.Close()
	}


	file1.Write([]byte(result.ExecuteInfo.Job.Name+"\n"+string(result.Output)))
	 //fmt.Println("任务执行完成:", result.ExecuteInfo.Job.Name, string(result.Output), result.Err)
}

// 调度协程
func (scheduler *Scheduler) scheduleLoop() {
	var (
		jobEvent *job.JobEvent
		scheduleAfter time.Duration
		scheduleTimer *time.Timer
		jobResult *job.JobExecuteResult
	)

	// 初始化一次(1秒)
	scheduleAfter = scheduler.TrySchedule()

	// 调度的延迟定时器
	scheduleTimer = time.NewTimer(scheduleAfter)

	// 定时任务common.Job
	for {
		select {
		case jobEvent = <- scheduler.jobEventChan:	//监听任务变化事件
			// 对内存中维护的任务列表做增删改查
			scheduler.handleJobEvent(jobEvent)
		case <- scheduleTimer.C:	// 最近的任务到期了
		case jobResult = <- scheduler.jobResultChan: // 监听任务执行结果
			scheduler.handleJobResult(jobResult)
		}
		// 调度一次任务
		scheduleAfter = scheduler.TrySchedule()
		// 重置调度间隔
		scheduleTimer.Reset(scheduleAfter)
	}
}

// 推送任务变化事件
func (scheduler *Scheduler) PushJobEvent(jobEvent *job.JobEvent) {
	scheduler.jobEventChan <- jobEvent
}

// 初始化调度器
func InitScheduler() (err error) {
	G_scheduler = &Scheduler{
		jobEventChan: make(chan *job.JobEvent, 1000),
		jobPlanTable: make(map[string]*job.JobSchedulePlan),
		jobExecutingTable: make(map[string]*job.JobExecuteInfo),
		jobResultChan: make(chan *job.JobExecuteResult, 1000),
	}
	// 启动调度协程
	go G_scheduler.scheduleLoop()
	return
}

// 回传任务执行结果
func (scheduler *Scheduler) PushJobResult(jobResult *job.JobExecuteResult) {
	scheduler.jobResultChan <- jobResult
}
