package request

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/sirupsen/logrus"
)

type base struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

type AccessToken struct {
	base
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}
type H map[string]any

func request(url string, method string, data io.Reader, header H) ([]byte, error) {
	req, err := http.NewRequest(method, url, data)
	if err != nil {
		return nil, err
	}
	for k, v := range header {
		req.Header.Set(k, v.(string))
	}
	req.Header.Set("Content-Type", "application/json")
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var tmp struct {
		url    string
		method string
		buf    string
		err    error
		param  string
	}
	tmp.url = url
	tmp.method = method
	tmp.buf = string(buf)
	tmp.err = err
	if data != nil {
		buf, err := io.ReadAll(data)
		if err == nil {
			tmp.param = string(buf)
		}
	}
	logrus.Warnf("request=>%+v", tmp)
	return buf, nil
}

func Get(url string, header H) ([]byte, error) {
	if data, err := request(url, "GET", nil, header); err != nil {
		logrus.Errorf("get %s, err:%s", url, err)
		return nil, err
	} else {
		return data, nil
	}
}
func Post(url string, data any, header H) ([]byte, error) {
	buf, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	reader := bytes.NewBuffer(buf)
	if data, err := request(url, "POST", reader, header); err != nil {
		logrus.Errorf("post %s, err:%s", url, err)
		return nil, err
	} else {
		return data, nil
	}
}
