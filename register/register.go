package register

import "github.com/qingcloud-iot/device-sdk-go/internal/utils/httpclient"

const (
	REGISTER_API = "/api/register/devices"
)

type Register struct {
	ServiceAddress string
}

func NewRegister(addr string) *Register {
	return &Register{
		ServiceAddress: addr,
	}
}

// DynamicRegistry 大批量设备的动态注册
func (t *Register) DynamicRegistry(midCredential string) (*httpclient.Data, error) {
	url := "http://" + t.ServiceAddress + REGISTER_API + "/" + midCredential
	client := httpclient.NewHttpClient(url)
	data, err := client.Post()
	if err != nil {
		return nil, err
	}
	return data, nil
}
