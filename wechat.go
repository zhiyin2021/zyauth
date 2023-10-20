package zyauth

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/zhiyin2021/zyauth/request"

	"github.com/sirupsen/logrus"
)

type wechatAuth struct {
	*authBase
}

func (w *wechatAuth) url(action string, args ...any) string {
	return fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/"+action, args...)
}
func (w *wechatAuth) GetAccessToken() string {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	if time.Since(w.access_last) < time.Second*7200 {
		return w.access_token
	}
	url := w.url(`gettoken?corpid=%s&corpsecret=%s`, w.appId, w.appSecret)
	buf, err := request.Get(url)
	if err != nil {
		logrus.Errorf("get accesstoken err:%s", err)
		return ""
	}
	var data map[string]any
	if err := json.Unmarshal(buf, &data); err != nil {
		logrus.Errorf("get accesstoken err:%s", err)
		return ""
	}
	w.access_token = data["access_token"].(string)
	w.access_last = time.Now().Add(time.Second * 7200)
	return w.access_token
}
func (w *wechatAuth) GetUserInfo(code string) *request.UserInfo {
	url := w.url(`user/getuserinfo?access_token=%s&code=%s`, w.GetAccessToken(), code)
	buf, err := request.Get(url)
	if err != nil {
		logrus.Errorf("get accesstoken err:%s", err)
		return nil
	}
	var data *request.UserInfo
	if err := json.Unmarshal(buf, &data); err != nil {
		logrus.Errorf("get accesstoken err:%s", err)
		return nil
	}
	return data
}
func (w *wechatAuth) GetUserDetail(userId string) *request.UserDetail {
	url := w.url(`user/get?access_token=%s&userid=%s`, w.GetAccessToken(), userId)
	buf, err := request.Get(url)
	if err != nil {
		logrus.Errorf("get accesstoken err:%s", err)
		return nil
	}
	var data *request.UserDetail
	if err := json.Unmarshal(buf, &data); err != nil {
		logrus.Errorf("get accesstoken err:%s", err)
		return nil
	}
	return data
}

func (w *wechatAuth) SendCard(toUser string, card request.WxMessageCard) error {
	url := w.url(`%smessage/send?access_token=%s`, w.GetAccessToken())
	wxmsg := request.WxMessage{
		AgentId: w.agentId,
		ToUser:  toUser,
		MsgType: "textcard",
		TextCard: request.WxMessageCard{
			Title:       card.Title,
			Description: card.ToDescription(),
			Url:         card.Url,
			Btntxt:      card.Btntxt,
		},
	}
	buf, err := json.Marshal(wxmsg)
	if err != nil {
		logrus.Errorf("send card err:%s", err)
		return err
	}
	if _, err := request.Post(url, buf); err != nil {
		logrus.Errorf("send card err:%s", err)
		return err
	}
	return nil
}
func (w *wechatAuth) SendMessage(toUser, content string) error {
	action := w.url(`%smessage/send?access_token=%s`, w.GetAccessToken())
	wxmsg := request.WxMessage{
		AgentId: w.agentId,
		ToUser:  toUser,
		MsgType: "text",
		Text: request.WxMessageText{
			Content: content,
		},
	}
	buf, err := json.Marshal(wxmsg)
	if err != nil {
		logrus.Errorf("send card err:%s", err)
		return err
	}
	if _, err := request.Post(action, buf); err != nil {
		logrus.Errorf("send card err:%s", err)
		return err
	}
	return nil
}

func (w *wechatAuth) AuthUrl(redirectUrl, state string) string {
	return fmt.Sprintf(`https://open.work.weixin.qq.com/wwopen/sso/qrConnect?appid=%s&agentid=%s&state=%s&redirect_uri=%s&self_redirect=true`,
		w.appId, w.agentId, state, url.QueryEscape(redirectUrl))
}

func (d *wechatAuth) loadData() {
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
func (d *wechatAuth) saveData() {
	data := fmt.Sprintf(`{"token":"%s","last":%d"`, d.access_token, d.access_last.Unix())
	os.WriteFile("zyauth.hisory", []byte(data), 0666)
}
