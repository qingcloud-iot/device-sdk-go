package mqtt

import (
	"context"
	"git.internal.yunify.com/tools/device-sdk-go/index"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

/**
* @Author: hexing
* @Date: 19-9-9 下午12:19
 */
func TestNewMqtt(t *testing.T) {
	options := &index.Options{
		Token:    "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY3IiOiIxIiwiYXVkIjoiaWFtIiwiYXpwIjoiaWFtIiwiY3VpZCI6ImlhbXItbW4xN3prNXoiLCJlaXNrIjoiTXk4YTBNRlZjRE93UG9Mc3lWOTlkSlVKU3QydmFrb1dmSTkxT21JeTc2dz0iLCJleHAiOjE1OTg5NDMyMzEsImlhdCI6MTU2NzQwNzIzMSwiaXNzIjoic3RzIiwianRpIjoiQTVMVjAzT1Bsc1ZuYndZa1R4Z2ZhUSIsIm5iZiI6MCwib3JnaSI6ImlvdGQtZDM3NTU2YjYtM2E1MS00ZTdhLWFkOTMtMmM4MjJjZTNmZWZkIiwib3d1ciI6InVzci1rTFZWQkRxZCIsInByZWYiOiJxcm46cWluZ2Nsb3VkOmlhbToiLCJydHlwIjoicm9sZSIsInN1YiI6InN0cyIsInRoaWQiOiI4IiwidHlwIjoiSUQifQ.G8HgdDM4Sk0Y0Ikf_OCHm2Gw3RX0HROWm8coeJV6jVMJCcPpx4_SgSBkXs7tlzoVJV6HTtx_0CeKw6yISx_ubTErhaHMww_6TyOQx0pzpzCKbHd2HaNYxyjCJJ34gYEmZZRtfuXlX1atdNmvOdgknSMyx8xVFV4v1eyv0U-UQXRrNIKOmV03lXGGEWx2NqV1GN8JjI4taZHrj1iVXLgmXBnEV_PoV67iqrtwoOeVqi7CcXx9GGs5hdL5Dca3lK8lsInC3xu80hqM0jnA5yXvSWOE4yRsvps3CkHCmaYxbHKUkkQRfDkEvljQQ_QLw0Yo_bmD-xMlAK8LPpN-GXgS7Q",
		DeviceId: "iotd-d37556b6-3a51-4e7a-ad93-2c822ce3fefd",
		Server:   "tcp://192.168.14.120",
	}
	m, err := NewMqtt(options)
	assert.Nil(t, err)
	assert.NotNil(t, m)
	err = m.Start(nil, nil)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	data := index.Metadata{
		"key1": "xxx",
		"key2": "yyy",
	}
	reply := m.PubProperty(ctx, data)
	t.Log(reply)
	assert.Nil(t, err)
}
