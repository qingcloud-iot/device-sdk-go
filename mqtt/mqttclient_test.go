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
		Token:    "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY3IiOiIxIiwiYXVkIjoiaWFtIiwiYXpwIjoiaWFtIiwiY3VpZCI6ImlhbXItMnE5MjMzbjEiLCJlaXNrIjoiNGNUT0luSFkzX2VURDNfblJmZDQwdl9LVVJFNnBNcEJDZ21uX3d3OTh5VT0iLCJleHAiOjE1OTg5NDMyMDYsImlhdCI6MTU2NzQwNzIwNiwiaXNzIjoic3RzIiwianRpIjoiQTVMVjAzT1Bsc1ZuYndZa1R4Z2ZaciIsIm5iZiI6MCwib3JnaSI6ImlvdGQtYjk5MWQ2NDAtODM1Ni00NDA4LTllMmQtYzA0MjJmODYxMjMzIiwib3d1ciI6InVzci1rTFZWQkRxZCIsInByZWYiOiJxcm46cWluZ2Nsb3VkOmlhbToiLCJydHlwIjoicm9sZSIsInN1YiI6InN0cyIsInRoaWQiOiJpb3R0LTY1YTMyYWZiLTBkMTUtNDQwYS1iNzBkLThhNTFkMWM1Y2NkZSIsInR5cCI6IklEIn0.wBO04t88gOWi4dzTCbFB-KiipuRXKKDjc9v8x2vFcStFTnHbbB4At8dHrKP-FtA006xfn6ORYvUBkQFd3HE_Oyn4DevWkEz1rrk9DNtlRm2U59ppUoBTi2OQJSTqs277gBnb5ApyI6VbHVAfcUf8tP6m9EBHvlQnwsIFMIzRnD8ouYuynZH5MVkqDHxEu36hFzjz8aareOPLupEgBEEi-nZPpyfAviJYue_8M273Ho7CdBV9akbF-pW6uNEjUF2ejuan_dw9terACbjLAoGYOckxqIfQFqD9yiwnQ3AymVwh9zx4aXOT0neWk1VnH9FZCsZT0gp4Jpm4ocJ0L10cqg",
		DeviceId: "iotd-b991d640-8356-4408-9e2d-c0422f861233",
		Server:   "tcp://192.168.14.120:8055",
	}
	m, err := NewMqtt(options)
	assert.Nil(t, err)
	assert.NotNil(t, m)
	err = m.Start(nil, nil)
	assert.Nil(t, err)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	data := index.Metadata{
		"key1": 111,
		"key2": "xxx",
	}
	reply, err := m.PubProperty(ctx, data)
	assert.Nil(t, err)
	t.Log(reply)
	//name := "test"
	//reply = m.PubEvent(ctx, name, data)
	//t.Log(reply)
}
