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
		Token: "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY3IiOiIxIiwiYXVkIjoiaWFtIiwiYXpwIjoiaWFtIiwiY3VpZCI6ImlhbXItbDZjdGJoZjAiLCJlaXNrIjoiR2VuazUxbm5BLXZyOUJaSnJQQ1gwNnNPSnBabElFZmw4eGlkVEFNbWRjQT0iLCJleHAiOjE2MDU5Njg2MTQsImlhdCI6MTU3NDQzMjYxNCwiaXNzIjoic3RzIiwianRpIjoiVWpJNFdQQW9wNWdQNldPdHJIUTc5USIsIm5iZiI6MCwib3JnaSI6ImlvdGQtZDhjYmEzOTItYWU0NC00MGRmLTk2YzgtNmQ3MWMzMmI4NjZlIiwib3d1ciI6InVzci16eUt5UFNmRyIsInByZWYiOiJxcm46cWluZ2Nsb3VkOmlhbToiLCJydHlwIjoicm9sZSIsInN1YiI6InN0cyIsInRoaWQiOiJpb3R0LVdDcnQ5bk1hUFMiLCJ0eXAiOiJJRCJ9.V0hqewKk6cwwlWzUpBY1HFpMcEvElurmKHh_HtAD816oVsEvl58kK4zpfs1jslASfBLw11OHBE-BD1Zp9FfGicRgTulQ2OUI4t9UiDbmnxGGKODknuP-0lEAb30n6JqLWWZh-rlZlN0tQVixelMC45ftf4LR0OmRH1T250RWO1MNNqqNgral9juTZ8mI9qcvX0yN3Ro7hM_JndeFWc4j9uj_QLus-Sv0mhleMh4i_5uoji7p8XReykwC82Lm2o61EGZZ3T7RCW9GCrSFngIsXnFUxk9mGqUiyW4aqKNkvpcCg-lm3t4fuszc6YW9_YzU53uic14ERRswREf3Wj3vJg",
		//DeviceId: "iotd-78c675f2-3495-4734-af5b-bb31f83f764c",
		Server: "tcp://172.31.141.199:1889",
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

func TestTask(t *testing.T) {
	cache := cache.Cache("xxxx")
	cache.Add("a", 3*time.Second, "xxx")
	cache.RemoveAboutToDeleteItemCallback()
}
