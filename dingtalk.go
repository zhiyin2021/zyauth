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

type dingtalkAuth struct {
	*authBase
}

// {HTTP method} https://api.dingtalk.com/{version}/{resource}?{query-parameters}
func (d *dingtalkAuth) api(action string, args ...any) string {
	// return fmt.Sprintf("https://oapi.dingtalk.com/"+action, args...)
	return fmt.Sprintf("https://api.dingtalk.com/"+action, args...)
}
func (d *dingtalkAuth) topapi(action string, args ...any) string {
	// return fmt.Sprintf("https://oapi.dingtalk.com/"+action, args...)
	return fmt.Sprintf("https://oapi.dingtalk.com/"+action, args...)
}
func (d *dingtalkAuth) GetAccessToken() string {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	if time.Since(d.access_last) < time.Second*7200 {
		return d.access_token
	}
	url := d.topapi(`gettoken?appkey=%s&appsecret=%s`, d.appId, d.appSecret)
	buf, err := request.Get(url, nil)
	if err != nil {
		logrus.Errorf("get accesstoken err:%s", err)
		return ""
	}
	var data request.AccessToken
	if err := json.Unmarshal(buf, &data); err != nil {
		logrus.Errorf("get accesstoken err:%s", err)
		return ""
	}
	d.access_token = data.AccessToken
	d.access_last = time.Now().Add(time.Second * 7200)
	d.saveData()
	return d.access_token
}
func (d *dingtalkAuth) getUserToken(code string) string {
	url := d.api(`v1.0/oauth2/userAccessToken`)
	// https://api.dingtalk.com/v1.0/oauth2/userAccessToken method:POST buf:{"expireIn":7200,"accessToken":"f992392321953e9892076b0c748b75cb","refreshToken":"70516328885f3edc8aa491a961178859"} err:<nil> param:}
	buf, err := request.Post(url, request.H{
		"clientId":     d.appId,
		"clientSecret": d.appSecret,
		"code":         code,
		"grantType":    "authorization_code",
	}, nil)
	if err != nil {
		logrus.Errorf("getUserToken err:%s", err)
		return ""
	}
	var data request.H
	if err := json.Unmarshal(buf, &data); err != nil {
		logrus.Errorf("getUserToken json.unmarshal[%s] err:%s", buf, err)
		return ""
	}
	if v, ok := data["accessToken"].(string); ok {
		return v
	}
	return ""

	// https://api.dingtalk.com/v1.0/oauth2/userAccessToken method:POST buf:{"code":"Missingbody","requestid":"90863DAE-1842-76EF-92F2-1ADEC7C89D6B","message":"body is mandatory for this action."} err:<nil> param:}

}

func (d *dingtalkAuth) GetUserInfo(code string) *request.UserInfo {
	action := d.api("v1.0/contact/users/me")
	buf, err := request.Get(action, request.H{"x-acs-dingtalk-access-token": d.getUserToken(code)})
	// https://api.dingtalk.com/v1.0/contact/users/me method:GET buf:{"nick":"金华","unionId":"vQR2XRTIXbAUYtz7oy0nfwiEiE","openId":"Uq1wkz7qyLzomuUviPEkPugiEiE","mobile":"19942422224","stateCode":"86"} err:<nil> param:}

	if err != nil {
		logrus.Errorf("GetUserInfo err:%s", err)
		return nil
	}
	var data request.DUserInfo
	if err := json.Unmarshal(buf, &data); err != nil {
		logrus.Errorf("GetUserInfo unjson %s, err:%s", buf, err)
		return nil
	}
	userInfo := &request.UserInfo{
		UserId:  data.Unionid,
		Unionid: data.Unionid,
		OpenId:  data.OpenId,
		Nick:    data.Nick,
		Mobile:  data.Mobile,
	}
	return userInfo
}

func (d *dingtalkAuth) GetUserDetail(unionId string) *request.UserDetail {
	userId := d.getUseridByUnionid(unionId)
	if userId == "" {
		return nil
	}
	url := d.topapi(`topapi/v2/user/get?access_token=%s`, d.GetAccessToken())
	// https://oapi.dingtalk.com/topapi/v2/user/get?access_token=40b36df674dc38209553ce5a939c4624 method:POST
	//buf:{"errcode":0,"errmsg":"ok","result":{"active":true,"admin":true,"avatar":"","boss":true,"dept_id_list":[1],"dept_order_list":[{"dept_id":1,"order":176220178909140512}],"exclusive_account":false,"hide_mobile":false,"leader_in_dept":[{"dept_id":1,"leader":false}],"mobile":"19942422224","name":"金华","real_authed":true,"role_list":[{"group_name":"默认","id":3935025514,"name":"主管理员"},{"group_name":"默认","id":3935025516,"name":"负责人"}],"senior":false,"state_code":"86","unionid":"vQR2XRTIXbAUYtz7oy0nfwiEiE","userid":"18293148161178525"},"request_id":"16lzmj9nqoklx"} err:<nil> param:}
	//get accesstoken err:json: cannot unmarshal array into Go struct field .result.dept_id_list of type string
	buf, err := request.Post(url, request.H{"userid": userId}, nil)
	if err != nil {
		logrus.Errorf("get accesstoken err:%s", err)
		return nil
	}
	var data request.DUserDetail
	if err := json.Unmarshal(buf, &data); err != nil {
		logrus.Errorf("get accesstoken err:%s", err)
		return nil
	}
	userDetail := &request.UserDetail{
		UserId:   data.Result.UserId,
		Name:     data.Result.Name,
		Mobile:   data.Result.Mobile,
		Avatar:   data.Result.Avatar,
		Position: data.Result.Title,
		Status:   1,
	}
	if data.Result.Active {
		userDetail.Status = 2
	}
	if data.Result.Active {
		userDetail.Status = 1
	}
	return userDetail
}

func (d *dingtalkAuth) SendCard(toUser string, card request.WxMessageCard) error {
	action := d.api(`topapi/message/send?access_token=%s`, d.GetAccessToken())
	wxmsg := request.WxMessage{
		AgentId: d.agentId,
		ToUser:  toUser,
		MsgType: "textcard",
		TextCard: request.WxMessageCard{
			Title:       card.Title,
			Description: card.ToDescription(),
			Url:         card.Url,
			Btntxt:      card.Btntxt,
		},
	}
	if _, err := request.Post(action, wxmsg, nil); err != nil {
		logrus.Errorf("send card err:%s", err)
		return err
	}
	return nil
}
func (d *dingtalkAuth) SendMessage(toUser, content string) error {
	action := d.api(`topapi/message/send?access_token=%s`, d.GetAccessToken())
	wxmsg := request.WxMessage{
		AgentId: d.agentId,
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
	if _, err := request.Post(action, buf, nil); err != nil {
		logrus.Errorf("send card err:%s", err)
		return err
	}
	return nil
}

func (d *dingtalkAuth) AuthUrl(redirectUrl, state string) string {
	return fmt.Sprintf(`https://login.dingtalk.com/oauth2/auth?response_type=code&client_id=%s&scope=openid&state=%s&redirect_uri=%s&prompt=consent`,
		d.appId, state, url.QueryEscape(redirectUrl))
}

func (d *dingtalkAuth) loadData() {
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
func (d *dingtalkAuth) saveData() {
	data := fmt.Sprintf(`{"token":"%s","last":%d}`, d.access_token, d.access_last.Unix())
	os.WriteFile("zyauth.hisory", []byte(data), 0666)
}

func (d *dingtalkAuth) getUseridByUnionid(unionId string) string {
	url := d.topapi("topapi/user/getbyunionid?access_token=%s", d.GetAccessToken())
	buf, err := request.Post(url, request.H{"unionid": unionId}, nil)
	if err != nil {
		logrus.Errorf("getUseridByUnionid request.err:%s", err)
		return ""
	}
	var data request.DUserId
	if err := json.Unmarshal(buf, &data); err != nil {
		logrus.Errorf("getUseridByUnionid unjson.err:%s", err)
		return ""
	}
	if data.ErrCode == 0 {
		return data.Result.UserId
	}
	return ""
}
