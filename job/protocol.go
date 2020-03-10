package job

import (
	"encoding/json"
	"strings"
	"../utils"
)


func GetJObName(oldname string)string  {
	s:=strings.TrimPrefix(oldname,utils.JOB_SAVE_DIR)
	return s
}
// 从 /cron/killer/job10提取job10
func ExtractKillerName(killerKey string) (string) {
	return strings.TrimPrefix(killerKey, utils.JOB_KILLER_DIR)
}
// 反序列化Job
func UnpackJob(value []byte) (ret *Job, err error) {
	var (
		job1 *Job
	)

	job1 = &Job{}
	if err = json.Unmarshal(value, job1); err != nil {
		return
	}
	ret = job1
	return
}
// 从etcd的key中提取任务名
// /cron/jobs/job10抹掉/cron/jobs/
func ExtractJobName(jobKey string) (string) {
	return strings.TrimPrefix(jobKey, utils.JOB_SAVE_DIR)
}
// 提取worker的IP
func ExtractWorkerIP(regKey string) (string) {
	return strings.TrimPrefix(regKey, utils.JOB_WORKER_DIR)
}
