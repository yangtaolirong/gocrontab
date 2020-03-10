package utils

import (
	//"fmt"
	"github.com/astaxie/beego/config"
	"log"
)
var Config config.Configer
func init()  {
	cf, err := config.NewConfig("ini", "./config.conf")
	if err != nil {
		log.Fatal("config error!",err)
		return
	}
	Config=cf
}
