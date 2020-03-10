package utils

import (
	"encoding/json"
	"log"
)

type Response struct {
	Errorcode int         `json:"code"`
	ErrMsg    string      `json:"errmsg"`
	Data      interface{} `json:"data"`
}

func BuildRes(errcode int,errmsg string,data interface{})[]byte {
	res:=Response{Errorcode:errcode,ErrMsg:errmsg, Data:data}
	resdata,err:=json.Marshal(&res)
	if err!=nil{
		log.Fatal(err)
	}
	return resdata

}
