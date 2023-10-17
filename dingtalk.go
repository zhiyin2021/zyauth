package zyauth

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/zhiyin2021/zyauth/request"

	"github.com/sirupsen/logrus"
)

type dingtalkAuth struct {
	*authBase
}

func (d *dingtalkAuth) url(action string, args ...any) string {
	return fmt.Sprintf("https://oapi.dingtalk.com/"+action, args...)
}
func (d *dingtalkAuth) GetAccessToken() string {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	if time.Since(d.access_last) < time.Second*7200 {
		return d.access_token
	}
	url := d.url(`gettoken?appkey=%s&appsecret=%s`, d.appId, d.appSecret)
	buf, err := request.Get(url)
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
	return d.access_token
}
func (d *dingtalkAuth) GetUserInfo(code string) *request.UserInfo {
	url := d.url(`topapi/v2/user/getuserinfo?access_token=%s`, d.GetAccessToken())
	buf, err := request.Post(url, []byte(fmt.Sprintf(`{"code":"%s"}`, code)))
	if err != nil {
		logrus.Errorf("get accesstoken err:%s", err)
		return nil
	}
	var data *request.DUserInfo
	if err := json.Unmarshal(buf, &data); err != nil {
		logrus.Errorf("get accesstoken err:%s", err)
		return nil
	}
	userInfo := &request.UserInfo{
		UserId: data.Result.Userid,
	}
	return userInfo
}

func (d *dingtalkAuth) GetUserDetail(userId string) *request.UserDetail {
	url := d.url(`topapi/v2/user/get?access_token=%s`, d.GetAccessToken())
	buf, err := request.Post(url, []byte(fmt.Sprintf(`{"userid":"%s"}`, userId)))
	if err != nil {
		logrus.Errorf("get accesstoken err:%s", err)
		return nil
	}
	var data *request.DUserDetail
	if err := json.Unmarshal(buf, &data); err != nil {
		logrus.Errorf("get accesstoken err:%s", err)
		return nil
	}
	userDetail := &request.UserDetail{
		UserId: data.Result.Userid,
		Name:   data.Result.Name,
		Mobile: data.Result.Mobile,
		Avatar: data.Result.Avatar,
		Status: 0,
	}
	if data.Result.Active {
		userDetail.Status = 1
	}
	return userDetail
}

func (d *dingtalkAuth) SendCard(toUser string, card request.WxMessageCard) error {
	action := d.url(`message/send?access_token=%s`, d.GetAccessToken())
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
func (d *dingtalkAuth) SendMessage(toUser, content string) error {
	action := d.url(`message/send?access_token=%s`, d.GetAccessToken())
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
	if _, err := request.Post(action, buf); err != nil {
		logrus.Errorf("send card err:%s", err)
		return err
	}
	return nil
}

func (d *dingtalkAuth) AuthUrl(redirectUrl, state string) string {
	return fmt.Sprintf(`https://login.dingtalk.com/oauth2/auth?response_type=code&client_id=%s&scope=openid&state=%s&redirect_uri=%s&prompt=consent`,
		d.appId, state, url.QueryEscape(redirectUrl))
}

// func OnChatReceive(ctx context.Context, data *chatbot.BotCallbackDataModel) ([]byte, error) {
// 	return []byte("success"), nil
// }

// func Init() {
// 	logger.SetLogger(logger.NewStdTestLogger())
// 	cli := client.NewStreamClient(
// 		client.WithAppCredential(client.NewAppCredentialConfig(env.Config().GetString("dingtalk.appId"), env.Config().GetString("dingtalk.appSecret"))),
// 		client.WithUserAgent(client.NewDingtalkGoSDKUserAgent()),
// 		client.WithSubscription(utils.SubscriptionTypeKCallback, "", chatbot.NewDefaultChatBotFrameHandler(OnChatReceive).OnEventReceived),
// 	)

// 	err := cli.Start(context.Background())
// 	if err != nil {
// 		panic(err)
// 	}

// 	defer cli.Close()

//		select {}
//	}
