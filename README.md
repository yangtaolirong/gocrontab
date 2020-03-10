# gocrontab
基于go语言iris框架开发 crontab web界面管理工具
使用方法：
1. 搭建etcd环境
2.修改config的etcd配置，ip修改为etcd 所在ip，端口指定httpport
3.go run main.go
4.浏览器访问https://domain:httpport
5.cron 表达式兼容linux crontab表达式，同时支持更精细，能达到秒级别，具体参考github.com/kataras/iris/context

