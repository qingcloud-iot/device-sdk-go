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
		Token: "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY3IiOiIxIiwiYXVkIjoiaWFtIiwiYXpwIjoiaWFtIiwiY3VpZCI6ImlhbXItbDZjcTh5emwiLCJlaXNrIjoiZGVaRGUwcG15YW9QTS1hd2R5VlIzS3JNbDNaa2FjR00xVnlkOUhzYjRIdz0iLCJleHAiOjE2MDQ3MjEwODEsImlhdCI6MTU3MzE4NTA4MSwiaXNzIjoic3RzIiwianRpIjoiVWpJNFdQQW9wNWdQNldPdHJIUTUwTSIsIm5iZiI6MCwib3JnaSI6ImlvdGQtNzhjNjc1ZjItMzQ5NS00NzM0LWFmNWItYmIzMWY4M2Y3NjRjIiwib3d1ciI6InVzci1rTFZWQkRxZCIsInByZWYiOiJxcm46cWluZ2Nsb3VkOmlhbToiLCJydHlwIjoicm9sZSIsInN1YiI6InN0cyIsInRoaWQiOiJpb3R0LVUxaElnSTBkMUEiLCJ0eXAiOiJJRCJ9.uHjxinx_BotQkRMCaL-DGtLwm4uLELwWE8Kv00sVxr92uEJEQ7NwRmtAlKj6IJnMcjQBkXJv-R4ceDt-fX2tJ0helyPhrCu-5EbU4G3FLQJLc3J4Cy8hd2Ltn0V3P_JolArSVWf4qQDGHYHI1nYXowh4D7J8ApACT2xQVcFiYzkET66Dmjczfsug9F318rUJyUxRqSBZS-7rVMr8dgyicjm6zhVThLlaq2xpHQWo6443szyeV2BoipWDFpA1dBDW4LgwXVb6nVs0ZxuPEaOsKWeciRxd2NC25Pfg9pvr6rvVlzL-2mA7_Budi5PYG_LeQ-LfA82LpJo8f--WgMTUtw",
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

	go func() {
		var i int64 = 1539362482000
		for {
			data := index.Metadata{
				"CO2Concentration": RandInt64(1, 100),
				"humidity":         RandInt64(1, 100),
			}
			//tm := (time.Now().Unix() - int64(i*60*60)) * 1000
			ch, err := m.PubPropertyAsyncEx(data, i)
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
			i = i + 60000
			if i > 1570898482000 {
				i = 1539362482000
			}
			time.Sleep(500 * time.Millisecond)
		}
		//os.Exit(0)
	}()
	go func() {
		var i int64 = 1570898482000
		for {
			data := index.Metadata{
				"CO2Concentration": RandInt64(1, 100),
				"humidity":         RandInt64(1, 100),
			}
			ch, err := m.PubPropertyAsyncEx(data, i)
			assert.Nil(t, err)
			select {
			case value := <-ch:
				fmt.Println(value)
			}
			i = i + 60000
			i = i + 60000
			if i > 1574093894607 {
				i = 1539362482000
			}
			time.Sleep(100 * time.Millisecond)
		}
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
