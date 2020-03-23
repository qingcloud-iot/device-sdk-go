package register

import "git.internal.yunify.com/iot-sdk/device-sdk-go/interval/httpclient"

const (
	REGISTER_URL = "http://192.168.14.121:8889/api/register/devices"
)

type Register struct {
}

func NewTaskInfo() *Register {
	return &Register{}
}

func (t *Register) DynamicRegistry(midCredential string) (*httpclient.Data, error) {
	url := REGISTER_URL + "/" + midCredential
	client := httpclient.NewHttpClient(url)
	data, err := client.Post()
	if err != nil {
		return nil, err
	}
	return data, nil
}
