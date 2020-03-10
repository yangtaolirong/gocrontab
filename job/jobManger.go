package job

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"go.etcd.io/etcd/clientv3"
	//``"go.etcd.io/etcd/lease"
	"time"
	"../utils"
)

type JobManager struct {
	client *clientv3.Client
	kv clientv3.KV
	lease clientv3.Lease
}
var JM *JobManager
func InitJobManager()  {
	JM=new(JobManager)
	timenum,err:=utils.Config.Int("etcd::DialTimeout")
	if err!=nil{
		fmt.Println(err)
		return
	}
	cf:=clientv3.Config{
		Endpoints:[]string{utils.Config.String("etcd::Endpoint")},
		DialTimeout:time.Duration(timenum)*time.Second}
	client,err:=clientv3.New(cf)
	if err!=nil{
		fmt.Println(err)
		return
	}
	lease:=clientv3.NewLease(client)//创建租约对象
	kv:=clientv3.NewKV(client)//创建kv对象
	JM.kv=kv
	JM.client=client
	JM.lease=lease

}
func (jm * JobManager)JobSave(job *Job)  (*Job,error) {
	key:=utils.JOB_SAVE_DIR+job.Name
	value,err:=json.Marshal(job)
	if err!=nil{
		return nil,err
	}
	fmt.Println(job)
	res,err:=jm.kv.Put(context.TODO(),key,string(value),clientv3.WithPrevKV())
	if err!=nil{
		fmt.Println(err)
		return nil, err
	}
	oldjob:=new(Job)
	if res.PrevKv!=nil{//获取之前的kv
		fmt.Println("res.prevkey=",res.PrevKv)
		err=json.Unmarshal(res.PrevKv.Value,&oldjob)
		if err!=nil{
			fmt.Println(err)
			return nil,err
		}
	}
	return oldjob,nil
}
func (jm * JobManager)JobDelete(key string)  (*Job,error){
	res,err:=jm.kv.Delete(context.TODO(),key)
	if err!=nil{
		fmt.Println(err)
		return nil,err
	}
	oldjob:=new(Job)
	if len(res.PrevKvs)!=0{
		err=json.Unmarshal(res.PrevKvs[0].Value,&oldjob)
		if err!=nil{
			return nil,err
		}
	}
	return oldjob,nil
}
func (jm * JobManager)JobList()  (jobs []Job,err error)  {
	res,err:=jm.kv.Get(context.TODO(),utils.JOB_SAVE_DIR,clientv3.WithPrefix())
	if err!=nil{
		fmt.Println(err)
		return nil,err
	}
	if len(res.Kvs)==0{
		return nil,errors.New("没有任务可以")
	}
	var job Job
	for _,v:=range res.Kvs{
		err:=json.Unmarshal(v.Value,&job)
		if err!=nil{
			return nil,err
		}
		jobs=append(jobs, job)
	}
	return jobs ,nil
}
func (jm * JobManager)JobKill(name string)  error {
	killdir := utils.JOB_KILLER_DIR + name
	l, err := jm.lease.Grant(context.TODO(), 1) //申请租约，时间为10秒
	if err != nil {
		return err
	}
	id := l.ID                                                                      //获取租约id
	_, err = jm.kv.Put(context.TODO(), killdir, "", clientv3.WithLease(id)) //关联租约
	if err!=nil {
		return err
	}
	return nil
}
