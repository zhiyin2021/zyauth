package request

// {
// 	"accessToken" : "abcd",
// 	"refreshToken" : "abcd",
// 	"expireIn" : 7200,
// 	"corpId" : "corpxxxx"
//   }

type DUserAccessToken struct {
	AccessToken  string `json:"accessToken"`
	ExpireIn     int    `json:"expireIn"`
	CorpId       string `json:"corpId"`
	RefreshToken string `json:"refreshToken"`
}
type DUserId struct {
	base
	Result struct {
		UserId string `json:"userid"`
	} `json:"result"`
}

type DUserInfo struct {
	Unionid string `json:"unionid"`
	OpenId  string `json:"openid"`
	Nick    string `json:"nick"`
	Mobile  string `json:"mobile"`
	Avatar  string `json:"avatarUrl"`
}

type DUserDetail struct {
	base
	Result struct {
		UserId           string `json:"userid"`            // 员工在当前企业内的唯一标识，也称staffId
		Unionid          string `json:"unionid"`           // 员工在当前开发者企业账号范围内的唯一标识
		Name             string `json:"name"`              // 员工姓名
		Boss             bool   `json:"boss"`              // 是否为企业的老板，true表示是，false表示不是
		Mobile           string `json:"mobile"`            // 手机号码
		Active           bool   `json:"active"`            // 是否已经激活，true表示已激活，false表示未激活
		ManagerUserid    string `json:"manager_userid"`    // 上级领导的userid
		Title            string `json:"title"`             // 职位信息
		ExclusiveAccount bool   `json:"exclusive_account"` // 是否为企业账号

		RoleList []struct {
			GroupName string `json:"group_name"`
			Name      string `json:"name"`
			Id        string `json:"id"`
		} `json:"role_list"`
		Admin     bool   `json:"admin"`
		Remark    string `json:"remark"`
		HiredDate int64  `json:"hired_date"`

		WorkPlace     string `json:"work_place"`
		DeptOrderList []struct {
			DeptId string `json:"dept_id"`
			Order  string `json:"order"`
		} `json:"dept_order_list"`
		RealAuthed bool `json:"real_authed"`
		// DeptIdList   string `json:"dept_id_list"`
		JobNumber    string `json:"job_number"`
		Email        string `json:"email"`
		LeaderInDept struct {
			Leader bool   `json:"leader"`
			DeptId string `json:"dept_id"`
		} `json:"leader_in_dept"`
		OrgEmail    string `json:"org_email"`
		Telephone   string `json:"telephone"`
		Avatar      string `json:"avatar"`
		HideMobile  bool   `json:"hide_mobile"`
		Senior      bool   `json:"senior"`
		UnionEmpExt struct {
			UnionEmpMapList []struct {
				Userid string `json:"userid"`
				CorpId string `json:"corp_id"`
			} `json:"union_emp_map_list"`
			Userid string `json:"userid"`
			CorpId string `json:"corp_id"`
		} `json:"union_emp_ext"`
		StateCode string `json:"state_code"`
	} `json:"result"`
}
