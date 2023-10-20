package zyauth

import (
	"sync"
	"time"

	"github.com/zhiyin2021/zyauth/request"
)

type Auth interface {
	GetAccessToken() string
	GetUserDetail(userId string) *request.UserDetail
	GetUserInfo(code string) *request.UserInfo
	AuthUrl(redirectUrl, state string) string
	saveData()
	loadData()
}
type authBase struct {
	access_token string
	access_last  time.Time
	mutex        sync.Mutex
	appId        string
	appSecret    string
	agentId      string
}

func NewAuth(tp, appId, appSecret, agentId string) Auth {
	base := &authBase{
		appId:     appId,
		appSecret: appSecret,
		agentId:   agentId,
	}
	var auth Auth
	switch tp {
	case "wechat":
		auth = &wechatAuth{base}
	case "dingtalk":
		auth = &dingtalkAuth{base}
	default:
		panic("not support auth type")
	}
	auth.loadData()
	return auth
}
