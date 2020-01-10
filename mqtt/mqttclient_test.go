package mqtt

import (
	"context"
	"fmt"
	"git.internal.yunify.com/tools/device-sdk-go/index"
	cache "github.com/muesli/cache2go"
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
		Token: "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY3IiOiIxIiwiYXVkIjoiaWFtIiwiYXpwIjoiaWFtIiwiY3VpZCI6ImlhbXItN2R4eWVvOTUiLCJlaXNrIjoieFd5SlRFcUR1MkozVGR1S3V1TVB6X1R4Z0lXZWZIVTRidVZTUG1yUlRQQT0iLCJleHAiOjE2MDkyMTQ0ODUsImlhdCI6MTU3NzY3ODQ4NSwiaXNzIjoic3RzIiwianRpIjoiVWpJNFdQQW9wNWdQNldPdHJIUUhXWiIsIm5iZiI6MCwib3JnaSI6ImlvdGQtY2JjNmQwZGYtOTBiYi00YzNhLTllM2MtMmY4MDNlNDVhNDczIiwib3d1ciI6InVzci1PN2xUUVY3NyIsInByZWYiOiJxcm46cWluZ2Nsb3VkOmlhbToiLCJydHlwIjoicm9sZSIsInN1YiI6InN0cyIsInRoaWQiOiJpb3R0LUZSRHA3cTZjSnUiLCJ0eXAiOiJJRCJ9.rFsWSeP6QIutmdyb13kBm9pqypMBz0Bt4uwHt-ER9kJpcPBy5mWN88yr7IXZ-_Xd1YcCr90NZYgl9bZILSM4J6S4WZ97g4W82ZMGIP8-mdq_dDV1bP-2bWREQ1upZKQdfi0fM7Q4-ZRWLx9Qsk7ZZl-7SFaJthmzXMEe9vTI0uYbZJ0T_FY8mEJybMhkht4r14NeuC4rr56Ci4paX6dBKFLYeKd_Jgijcb04-H4z_nvF1JUAABKuNV07LBkHX2hRLPu3P640T5hjNemP9Ya6D24P9uAKL4vJvABFFgAQkFOvB3W708U1CYMLZIfPVfZv0EVW1nutwhY_S4RP6Bj4Kg",
		//DeviceId: "iotd-78c675f2-3495-4734-af5b-bb31f83f764c",
		Server: "tcp://192.168.14.120:8055",
		SetProperty: func(meta index.Metadata) (index.Metadata, error) {
			fmt.Println("SetProperty", meta)
			data := make(index.Metadata)
			return data, nil
		},
		ServiceHandle: func(name string, meta index.Metadata) (index.Metadata, error) {
			fmt.Println("ServiceHandle", name, meta)
			data := make(index.Metadata)
			return data, nil
		},
	}
	m, err := NewMqtt(options)
	assert.Nil(t, err)
	assert.Nil(t, m)
	time.Sleep(5 * time.Second)
	//data := index.Metadata{
	//	"target":  "car",
	//	"color": "red",
	//	"license": "120",
	//	"address": "http://iot.qingcloud.com/test.jpg",
	//}
	data := index.Metadata{
		"target": "left",
		"cause":  "fail",
	}
	for {
		reply, err := m.PubEventSync(context.Background(), "FailureReport", data)
		assert.Nil(t, err)
		fmt.Println(reply)
		time.Sleep(1 * time.Second)
	}
	//go func() {
	//	var i int64 = 1539362482000
	//	for {
	//		data := index.Metadata{
	//			"CO2Concentration": RandInt64(1, 100),
	//			"humidity":         RandInt64(1, 100),
	//		}
	//		//tm := (time.Now().Unix() - int64(i*60*60)) * 1000
	//		ch, err := m.PubPropertyAsyncEx(data, i)
	//		//ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	//		//reply, err := m.PubPropertySync(ctx, data)
	//		//cancel()
	//		assert.Nil(t, err)
	//		select {
	//		case value := <-ch:
	//			fmt.Println(value)
	//		}
	//		//data = index.Metadata{
	//		//	"int32":  10,
	//		//	"string": "hexing-string",
	//		//	"float":  rand.Float32(),
	//		//	"double": rand.Float64(),
	//		//}
	//		//reply, err = m.PubEventSync(context.Background(), "he-event1", data)
	//		//assert.Nil(t, err)
	//		//fmt.Println(reply)
	//		i = i + 60000
	//		if i > 1570898482000 {
	//			i = 1539362482000
	//		}
	//		time.Sleep(10000 * time.Millisecond)
	//	}
	//	//os.Exit(0)
	//}()
	select {}
	//name := "test"
	//reply = m.PubEvent(ctx, name, data)
	//t.Log(reply)
}
func TestNewHubMqtt(t *testing.T) {
	options := &index.Options{
		SetProperty: func(meta index.Metadata) (index.Metadata, error) {
			fmt.Println("SetProperty", meta)
			data := make(index.Metadata)
			return data, nil
		},
		ServiceHandle: func(name string, meta index.Metadata) (index.Metadata, error) {
			fmt.Println("ServiceHandle", name, meta)
			data := make(index.Metadata)
			return data, nil
		},
	}
	m, err := NewHubMqtt(options)
	assert.Nil(t, err)
	assert.Nil(t, m)
	time.Sleep(5 * time.Second)
}
func TestTask(t *testing.T) {
	cache := cache.Cache("xxxx")
	cache.Add("a", 3*time.Second, "xxx")
	cache.RemoveAboutToDeleteItemCallback()
}
