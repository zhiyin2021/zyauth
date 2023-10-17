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
	switch tp {
	case "wechat":
		return &wechatAuth{base}
	case "dingtalk":
		return &dingtalkAuth{base}
	default:
		panic("not support auth type")
	}
}
