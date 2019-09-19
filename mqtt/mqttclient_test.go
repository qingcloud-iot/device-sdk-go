package mqtt

import (
	"context"
	"fmt"
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
		Token:    "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY3IiOiIxIiwiYXVkIjoiaWFtIiwiYXpwIjoiaWFtIiwiY3VpZCI6ImlhbXItOXcwZ3RkYW4iLCJlaXNrIjoiRTVVbzdzLVRkTkUybDJpeU14N2ZQc2t4WE9renR3akVTczliS3BGbXhVdz0iLCJleHAiOjE2MDAzMjc1NjgsImlhdCI6MTU2ODc5MTU2OCwiaXNzIjoic3RzIiwianRpIjoiQTVMVjAzT1Bsc1ZuYndZa1R4Z2ZqUyIsIm5iZiI6MCwib3JnaSI6ImlvdGQtNzUwNjQyODgtZGY4OC00M2JhLTgyMDctMTRmY2NhNjY5NzEyIiwib3d1ciI6InVzci1rTFZWQkRxZCIsInByZWYiOiJxcm46cWluZ2Nsb3VkOmlhbToiLCJydHlwIjoicm9sZSIsInN1YiI6InN0cyIsInRoaWQiOiJpb3R0LTlkNWJmMDQxLTZiOTctNGIzZi05ZWUwLTY0MzQwYjI0NTQzZCIsInR5cCI6IklEIn0.RIO5PFu9eR9x-0tU9iNJOSbP4UqlgITstoqJ8BSLYcl8IjRJFQ-s1g_QMo6OkAh0zW1BMJS3zSP613zgwzb5z9TRc7eAitcBLBjmgMkCnUqMCKy_EP4Z_fvFRVkYHuBsgyr3Mz-NTgaYPujOuvdjZh934bxpBm792wlFkJgBfW-p9mvediAHPepqwx1OxGDLxPPPBmvcdEkznZk5DdEoMvgmA9zDVf_-OpZ9smvbU-8gka09Ph1LBxnd9NPJBItDeQf3ZmIl9ePTGcWTlxWETbEEn9QnAtSMfdaVJgUYBfysBRssGO8qk2UZe08wdSHI0bNNFsWkUD2CJGqMsAMTlw",
		DeviceId: "iotd-75064288-df88-43ba-8207-14fcca669712",
		Server:   "tcp://192.168.14.120:8055",
	}
	m, err := NewMqtt(options)
	assert.Nil(t, err)
	assert.NotNil(t, m)
	assert.Nil(t, err)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	data := index.Metadata{
		"event1": 1,
		"event2": "xxx",
		"event3": 12.2,
		"event4": 12.2,
	}
	reply, err := m.PubEvent(ctx, "eventxx", data)
	//data := index.Metadata{
	//	"car":       "car",
	//	"car_res":   "car_res",
	//	"filepath": "filepath",
	//	"label": "xxx",
	//	"label_res":"label_res",
	//}
	//reply, err := m.PubProperty(ctx,data)
	assert.Nil(t, err)
	fmt.Println(reply)
	select {}
	//name := "test"
	//reply = m.PubEvent(ctx, name, data)
	//t.Log(reply)
}
