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
		Extension string `json:"extension"`
		Unionid   string `json:"unionid"`
		Boss      bool   `json:"boss"`
		RoleList  []struct {
			GroupName string `json:"group_name"`
			Name      string `json:"name"`
			Id        string `json:"id"`
		} `json:"role_list"`
		ExclusiveAccount bool   `json:"exclusive_account"`
		ManagerUserid    string `json:"manager_userid"`
		Admin            bool   `json:"admin"`
		Remark           string `json:"remark"`
		Title            string `json:"title"`
		HiredDate        int64  `json:"hired_date"`
		Userid           string `json:"userid"`
		WorkPlace        string `json:"work_place"`
		DeptOrderList    []struct {
			DeptId string `json:"dept_id"`
			Order  string `json:"order"`
		} `json:"dept_order_list"`
		RealAuthed   bool   `json:"real_authed"`
		DeptIdList   string `json:"dept_id_list"`
		JobNumber    string `json:"job_number"`
		Email        string `json:"email"`
		LeaderInDept struct {
			Leader bool   `json:"leader"`
			DeptId string `json:"dept_id"`
		} `json:"leader_in_dept"`
		Mobile      string `json:"mobile"`
		Active      bool   `json:"active"`
		OrgEmail    string `json:"org_email"`
		Telephone   string `json:"telephone"`
		Avatar      string `json:"avatar"`
		HideMobile  bool   `json:"hide_mobile"`
		Senior      bool   `json:"senior"`
		Name        string `json:"name"`
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
