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
		Token:    "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY3IiOiIxIiwiYXVkIjoiaWFtIiwiYXpwIjoiaWFtIiwiY3VpZCI6ImlhbXItdzNuMWM4ZzIiLCJlaXNrIjoiQkNmXzQzczc3VElzZ3pNcVJzNlM0V0VVZHlyUEUydlV3eWVIbkIyeTBoaz0iLCJleHAiOjE2MDAzMjg1NTQsImlhdCI6MTU2ODc5MjU1NCwiaXNzIjoic3RzIiwianRpIjoiQTVMVjAzT1Bsc1ZuYndZa1R4Z2ZrMSIsIm5iZiI6MCwib3JnaSI6ImlvdGQtOGI3NzM0NjctNjVkOC00Y2JhLThjYTgtYzZlNmZjNzNjZDUwIiwib3d1ciI6InVzci1rTFZWQkRxZCIsInByZWYiOiJxcm46cWluZ2Nsb3VkOmlhbToiLCJydHlwIjoicm9sZSIsInN1YiI6InN0cyIsInRoaWQiOiJpb3R0LTNhNDc3MGJmLTI4ZDYtNDg1Zi1hZDYwLTQzZGU4NzVkMzdiMCIsInR5cCI6IklEIn0.xfCUbiT8mbhFtMDimWpujpxR19bqB9HFcwQ7yHoR828xRsqjJ_eLXoTbpIKt99dlAWaVbaL1UgZxr608LH6aeizExVKmJxoIIx7UwqDdyUGWANZLdK1NuENa-T3qHQ51EhSGyXU2xcJirt11dn4sqSQzYy-QI0zZIxreQvFvowPLkHl8sI0ulY63yCsD_2ipywNFqc03XCLk6Ey4Vzzl9SesLskZYX1bOku3rDn9HFXr1XtQe6TKCOsxGOp91NyRwynLXiCfaFK12UHr_u5UuS5EpXzpAkMN1CEEibNX0M4e9EKsx4iY9EN6-4li1egbCUovR5QDD9bCVljYQTERNQ",
		DeviceId: "iotd-8b773467-65d8-4cba-8ca8-c6e6fc73cd50",
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
		"car":       "hexing",
		"car_res":   "xxx",
		"label_res": 1212,
	}
	reply, err := m.PubProperty(ctx, data)
	assert.Nil(t, err)
	fmt.Println(reply)
	select {}
	//name := "test"
	//reply = m.PubEvent(ctx, name, data)
	//t.Log(reply)
}
