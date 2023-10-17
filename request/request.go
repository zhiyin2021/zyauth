package request

import (
	"bytes"
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

func request(url string, method string, data io.Reader) ([]byte, error) {
	req, err := http.NewRequest(method, url, data)
	if err != nil {
		return nil, err
	}
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var body []byte
	_, err = resp.Body.Read(body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func Get(url string) ([]byte, error) {
	if data, err := request(url, "GET", nil); err != nil {
		logrus.Errorf("get %s, err:%s", url, err)
		return nil, err
	} else {
		return data, nil
	}
}
func Post(url string, buf []byte) ([]byte, error) {
	reader := bytes.NewReader(buf)
	if data, err := request(url, "POST", reader); err != nil {
		logrus.Errorf("get %s, err:%s", url, err)
		return nil, err
	} else {
		return data, nil
	}
}
