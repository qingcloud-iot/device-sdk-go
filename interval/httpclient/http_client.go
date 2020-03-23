package httpclient

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	DefaultRetryTimes = 3
	DefaultTimeout    = 5 * time.Second
)

type HttpClient struct {
	url    string
	client *http.Client
	retry  int
}

type RespData struct {
	Code string      `json:"code"`
	Data interface{} `json:"data"`
}

type Data struct {
	ID         string `json:"id"`
	DeviceName string `json:"device_name"`
	Token      string `json:"token"`
}

func NewHttpClient(url string) *HttpClient {
	client := &http.Client{}
	return &HttpClient{url: url, client: client}
}

func (h *HttpClient) SetRetryTimes(retry int) {
	h.retry = retry
}

func (h *HttpClient) SetTimeout(timeout time.Duration) {
	h.client.Timeout = timeout
}

func (c *HttpClient) Post() (*Data, error) {

	var (
		resp *http.Response
		err  error
		i    int
	)

	if c.retry == 0 {
		c.retry = DefaultRetryTimes
	}
	if c.client.Timeout == 0 {
		c.client.Timeout = DefaultTimeout
	}

	for i = 0; i < c.retry; i++ {
		resp, err = post(c)
		if err == nil {
			break
		}
	}
	if i == c.retry && err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	r, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}
	respData := &RespData{}
	err = json.Unmarshal(r, respData)
	if err != nil {
		return nil, err
	}
	if respData.Code != "0" {
		return nil, errors.New(respData.Data.(string))
	}

	d, err := json.Marshal(respData.Data)
	if err != nil {
		return nil, err
	}
	var data Data
	err = json.Unmarshal(d, &data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func post(c *HttpClient) (*http.Response, error) {
	req, err := http.NewRequest("POST", c.url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
