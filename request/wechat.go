package request

import (
	"fmt"
	"time"
)

type WxPhone struct {
	PhoneNumber     string    `json:"phoneNumber"`
	PurePhoneNumber string    `json:"purePhoneNumber"`
	CountryCode     string    `json:"countryCode"`
	Watermark       Watermark `json:"watermark"`
}

type Code2SessionResult struct {
	base
	OpenId     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionId    string `json:"unionid"`
}
type UserInfo struct {
	base
	Nick    string `json:"nick"`
	UserId  string `json:"userId"`
	OpenId  string `json:"openId"`
	Mobile  string `json:"mobile"`
	Unionid string `json:"unionid"`
}

type UserDetail struct {
	base
	UserId    string    `json:"userId"`    // 员工在当前企业内的唯一标识，也称staffId
	Watermark Watermark `json:"watermark"` // 数据水印
	OpenId    string    `json:"openId"`    // 员工在当前开发者企业账号范围内的唯一标识
	NickName  string    `json:"nickName"`  // 员工姓名
	Gender    int       `json:"gender"`    // 性别 0：未知、1：男、2：女
	City      string    `json:"city"`      // 城市，如杭州市
	Province  string    `json:"province"`  // 省份，如浙江省为Zhejiang
	Country   string    `json:"country"`   // 国家，如中国为CN
	Avatar    string    `json:"avatar"`    // 头像url
	UnionId   string    `json:"unionId"`   // 员工在当前开发者企业账号范围内的唯一标识
	Name      string    `json:"name"`      // 员工姓名
	Mobile    string    `json:"mobile"`    // 手机号码
	Position  string    `json:"position"`  // 职位
	Status    int       `json:"status"`    // 员工在企业内的状态 1:已激活 2:已禁用 4:未激活
}
type Watermark struct {
	Appid     string `json:"appid"`
	Timestamp int    `json:"timestamp"`
}

type WxMessageCard struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Url         string   `json:"url"`
	Btntxt      string   `json:"btntxt"`
	Data        []string `json:"-"`
}
type WxMessageText struct {
	Content string `json:"content"`
}

type WxMessage struct {
	ToUser                 string        `json:"touser"`
	ToParty                string        `json:"toparty"`
	ToTag                  string        `json:"totag"`
	MsgType                string        `json:"msgtype"`
	AgentId                string        `json:"agentid"`
	Text                   WxMessageText `json:"text"`
	TextCard               WxMessageCard `json:"textcard"`
	Safe                   int           `json:"safe"`
	EnableIdTrans          int           `json:"enable_id_trans"`
	EnableDuplicateCheck   int           `json:"enable_duplicate_check"`
	DuplicateCheckInterval int           `json:"duplicate_check_interval"`
}

func (card *WxMessageCard) ToDescription() string {
	res := ""
	if len(card.Data) > 0 {
		res = fmt.Sprintf(`<div class="gray">%s</div>`, time.Now().Format("2006-01-02 15:04"))
		for _, item := range card.Data {
			res += fmt.Sprintf(`<div class="normal">%s</div>`, item)
		}
		res += `<div class="highlight">请尽快处理,尽快处理,快处理,处理!</div>`
	}
	return res
}

// public class WxMessage
// {
// 	public string touser { get; set; } = "";
// 	public string toparty { get; set; } = "";
// 	public string totag { get; set; } = "";

// 	public string msgtype { get; set; } = "text";
// 	public string agentid { get; set; } = "";
// 	public WxMessageText text { get; set; }
// 	public WxMessageCard textcard { get; set; }
// 	public int safe { get; set; } = 0;
// 	public int enable_id_trans { get; set; } = 0;
// 	public int enable_duplicate_check { get; set; } = 1;
// 	public int duplicate_check_interval { get; set; } = 1800;
// 	//"textcard" : {
// 	//    "title" : "领奖通知",
// 	//    "description" : "<div class=\"gray\">2016年9月26日</div> <div class=\"normal\">恭喜你抽中iPhone 7一台，领奖码：xxxx</div><div class=\"highlight\">请于2016年10月10日前联系行政同事领取</div>",
// 	//    "url" : "URL",
// 	//  "btntxt":"更多"
// 	//},
// }
// public class WxMessageCard
//     {
//         public WxMessageCard()
//         {
//             Data = new List<string>();
//         }
//         public string title { get; set; }
//         private string _description = "";
//         public string description
//         {
//             get
//             {
//                 if (Data.Count > 0)
//                 {
//                     // gray normal  highlight
//                     DateTime tm = DateTime.Now;
//                     string res = $"<div class=\"gray\">{tm.ToString("yyyy-MM-dd HH:mm")}</div> ";
//                     foreach (var item in Data)
//                     {
//                         res += $"<div class=\"normal\">{item}</div>";

//                     }
//                     res += " <div class=\"highlight\">请尽快处理,尽快处理,快处理,处理!</div>";
//                     return res;
//                 }
//                 return _description;
//             }
//             set { _description = value; }
//         }
//         public string url { get; set; } = "";
//         public string btntxt { get; set; } = "more";
//         [JsonIgnore]
//         public List<string> Data { get; set; }
//     }
