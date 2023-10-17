package main

import (
	"fmt"

	"github.com/zhiyin2021/zyauth"
)

const (
	tp        = "wechat" // "dingtalk"
	appId     = "xxxxxxxx"
	appSecret = "xxxxxxxxxxxxxx"
	agentId   = "1000000"
)

func main() {
	auth := zyauth.NewAuth(tp, appId, appSecret, agentId)
	token := "xxxxxxxxxxxxxx"
	authUrl := auth.AuthUrl("http://127.0.0.1:8080/auth.login", token)
	fmt.Println(authUrl)
}
