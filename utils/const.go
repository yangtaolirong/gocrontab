package utils

// 任务保存目录
var JOB_SAVE_DIR = "/cron/jobs/"
var WORKER_IPS="/cron/ips/"

// 任务强杀目录
var JOB_KILLER_DIR = "/cron/killer/"

// 任务锁目录
var JOB_LOCK_DIR = "/cron/lock/"

// 服务注册目录
var JOB_WORKER_DIR = "/cron/workers/"

// 保存任务事件
var JOB_EVENT_SAVE = 1

// 删除任务事件
var JOB_EVENT_DELETE = 2

// 强杀任务事件
var JOB_EVENT_KILL = 3
