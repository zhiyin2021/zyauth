package zyauth

import (
	"encoding/json"
	"fmt"
	"os"
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

func (d *authBase) loadData() {
	tmp, err := os.ReadFile("zyauth.hisory")
	if err == nil {
		var data struct {
			Token string `json:"token"`
			Last  int64  `json:"last"`
		}
		if err := json.Unmarshal(tmp, &data); err == nil {
			d.access_token = data.Token
			d.access_last = time.Unix(data.Last, 0)
		}
	}
}
func (d *authBase) saveData() {
	data := fmt.Sprintf(`{"token":"%s","last":%d}`, d.access_token, d.access_last.Unix())
	os.WriteFile("zyauth.hisory", []byte(data), 0666)
}
