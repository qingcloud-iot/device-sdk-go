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
		Token:    "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY3IiOiIxIiwiYXVkIjoiaWFtIiwiYXpwIjoiaWFtIiwiY3VpZCI6ImlhbXItNmQ4cnlrbjUiLCJlaXNrIjoiQlpaQm84TTZObldKb3Jqb1JjUWJSUVhnRlMwY3pXOXpuOFBqemZROC00RT0iLCJleHAiOjE2MDQzOTg5MTAsImlhdCI6MTU3Mjg2MjkxMCwiaXNzIjoic3RzIiwianRpIjoiVWpJNFdQQW9wNWdQNldPdHJIUTNEVCIsIm5iZiI6MCwib3JnaSI6ImlvdGQtNmMxNjE3MGUtOTZlMy00OTcyLTlkYmQtMDFkNjBlYjMyNDUwIiwib3d1ciI6InVzci1rTFZWQkRxZCIsInByZWYiOiJxcm46cWluZ2Nsb3VkOmlhbToiLCJydHlwIjoicm9sZSIsInN1YiI6InN0cyIsInRoaWQiOiJpb3R0LWZUeXRjS1BWTlEiLCJ0eXAiOiJJRCJ9.O-ha2kQ2JbXaQQ2IGxlP1YdbPBFDkRpEZ5MUyK9AOva1wzO5XKPNunfW6eUx-ANOivodr2zkhg-7gA10uV1G-KFbZouHYnxjR7Fa142mHufRoUowVwSnmXWA_PD26guxv4ObvfNxPzdg2bMgqOxQ0LP56KnQ1MfCVj_Awl5Esdw8U6CUiUaJOPaZLvB6Ps8SzjtR_EipQBBbthSvLWDqTy8tnULLB-BNjyAeAvU00ePPnbqsIXoD25ddUUBJ6YJEMlt2FwI1lrqNbjA2w-xlEZPQXOHAkpF4Djl7QBGVdGdS7yfkf2PdxsQ24nwk05_B32YoPT3_URtAvBRgo8UuzA",
		DeviceId: "iotd-6c16170e-96e3-4972-9dbd-01d60eb32450",
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

	go func() {
		i := 0
		for {
			data := index.Metadata{
				"pro1": RandInt64(1, 100),
				"pro2": RandInt64(1, 100),
				"pro3": RandInt64(1, 100),
			}
			t := (time.Now().Unix() - int64(i*60*60)) * 1000
			ch, err := m.PubPropertyAsyncEx(data, t)
			//ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			//reply, err := m.PubPropertySync(ctx, data)
			//cancel()
			assert.Nil(t, err)
			select {
			case value := <-ch:
				fmt.Println(value)
			}
			//data = index.Metadata{
			//	"int32":  10,
			//	"string": "hexing-string",
			//	"float":  rand.Float32(),
			//	"double": rand.Float64(),
			//}
			//reply, err = m.PubEventSync(context.Background(), "he-event1", data)
			//assert.Nil(t, err)
			//fmt.Println(reply)
			i++
			time.Sleep(5 * time.Second)
		}
		//os.Exit(0)
	}()
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
