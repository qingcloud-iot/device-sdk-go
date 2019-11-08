package mqtt

import (
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
		Token:    "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY3IiOiIxIiwiYXVkIjoiaWFtIiwiYXpwIjoiaWFtIiwiY3VpZCI6ImlhbXItMGUwcGlheTAiLCJlaXNrIjoiQW5ZclVVSF9sOUI4ekc2QXR5cFAxY0JJOGp2dkZmRkw0bDBzVm50T3N2UT0iLCJleHAiOjE2MDQ0OTEzMDMsImlhdCI6MTU3Mjk1NTMwMywiaXNzIjoic3RzIiwianRpIjoiVWpJNFdQQW9wNWdQNldPdHJIUTRNUiIsIm5iZiI6MCwib3JnaSI6ImlvdGQtMzkwZjI4MWQtMjEwYi00NzQ4LTlmMmMtMGZiZjQxMGFjZTg2Iiwib3d1ciI6InVzci1rTFZWQkRxZCIsInByZWYiOiJxcm46cWluZ2Nsb3VkOmlhbToiLCJydHlwIjoicm9sZSIsInN1YiI6InN0cyIsInRoaWQiOiJpb3R0LTF5U0ZJdk5IVXEiLCJ0eXAiOiJJRCJ9.Uw7oZZYRwp1l_riQ7vVqc2Dx-CV0dZticZmnEHR8QbK-GaoTS4k5y9qoGBAZOOO_UetIiSCm3NRHtMblC67SKabaE01_McNIyoL1jWO0X1kwt_dLSTScc7IeVrMH5wgl7C_rsx-vPT2xupsO2prLbN9SRlQJa8lT72x_x9CgLnCibH4fYUYs-QutFv84XHTuj8T7PbBlcPdYmLw3axB75LgduWlmtzb5L962-NlxvctG2IFXeg4bONGN7zjWn_NRGlmHN_qyUgEyZ2F5OtWusZfPwvns_QKFS-DKDhX_x8I9qZx0EH0cnj2Nl0k2CpuQlNJqHPmeMWAgfe09VP63Jg",
		DeviceId: "iotd-390f281d-210b-4748-9f2c-0fbf410ace86",
		Server:   "tcp://192.168.14.120:8055",
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

	//go func() {
	//	for {
	//		i := 24*30 * (RandInt64(1,100))
	//		data := index.Metadata{
	//			"pro1": RandInt64(1, 100),
	//			"pro2": RandInt64(1, 100),
	//			"pro3": RandInt64(1, 100),
	//		}
	//		tm := (time.Now().Unix() - int64(i*60*60)) * 1000
	//		ch, err := m.PubPropertyAsyncEx(data, tm)
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
	//		time.Sleep(1 * time.Second)
	//	}
	//	//os.Exit(0)
	//}()
	select {}
	//name := "test"
	//reply = m.PubEvent(ctx, name, data)
	//t.Log(reply)
}

func TestTask(t *testing.T) {
	cache := cache.Cache("xxxx")
	cache.Add("a", 3*time.Second, "xxx")
	cache.RemoveAboutToDeleteItemCallback()
}
