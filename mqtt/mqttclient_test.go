package mqtt

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"git.internal.yunify.com/tools/device-sdk-go/index"
	cache "github.com/muesli/cache2go"
	"github.com/stretchr/testify/assert"
)

/**
* @Author: hexing
* @Date: 19-9-9 下午12:19
 */
func TestNewMqtt(t *testing.T) {
	options := &index.Options{
		Token: "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY3IiOiIxIiwiYXVkIjoiaWFtIiwiYXpwIjoiaWFtIiwiY3VpZCI6ImlhbXItaDJjMjl2OXIiLCJlaXNrIjoiVjRSd3NoNjRpcXJhSTVJTHlnZ2xHZFhnV3E1S1JGWWxFYnRwakkxZk9Raz0iLCJleHAiOjE2MDQ1NDc0ODYsImlhdCI6MTU3MzAxMTQ4NiwiaXNzIjoic3RzIiwianRpIjoiVWpJNFdQQW9wNWdQNldPdHJIUTRlQiIsIm5iZiI6MCwib3JnaSI6ImlvdGQtMjk5MDY2MDktNzNiYS00NzBkLWE2ZmQtMGUxYzE3MTkwZmQwIiwib3d1ciI6InVzci1rTFZWQkRxZCIsInByZWYiOiJxcm46cWluZ2Nsb3VkOmlhbToiLCJydHlwIjoicm9sZSIsInN1YiI6InN0cyIsInRoaWQiOiJpb3R0LWZUeXRjS1BWTlEiLCJ0eXAiOiJJRCJ9.C1oCwaviLAsmb42mDXX4mXw2h0ccXYV8Kd5mAGkCpxpOFM7Rd7lOL2kGMJpvv_I5caOTlSNiFwMe2L2eXiA_dsZPBEW08dmzghLZXpVABFG7KJOrxT5t6WBYzVCOezt4CynSXheIs0NjSMZ5VBTdiEjj8GIi5iAIWUaYrEeFOlj3IZPp7ddr82rkog9OIDnHDvyXDK2MruKAb7xZ2QZFa0Wg1GKixFUhfT0iU37pQZbsGAduj-kB9z4o_ZwtP8gFko6AkW8WuBzzXhs35cQty2vXJ3ohxKnXtoiwChNfIQmNr8Cc7VJmQTmrQPrgmK3uMnxi02SQXsF2vd0HmpA_7A",
		//DeviceId: "iotd-78c675f2-3495-4734-af5b-bb31f83f764c",
		Server: "tcp://127.0.0.1:1883", // 192.168.14.120:8055
		SetProperty: func(meta index.Metadata) (index.Metadata, error) {
			fmt.Println("SetProperty", meta)
			data := make(index.Metadata)
			return data, nil
		},
		ServiceHandle: func(name string, meta index.Metadata) (index.Metadata, error) {
			fmt.Println("ServiceHandle", name, meta)
			data := make(index.Metadata)
			data["test"] = "xxxx"
			return data, nil
		},
	}
	m, err := NewMqtt(options)
	assert.Nil(t, err)
	assert.Nil(t, m)
	time.Sleep(5 * time.Second)
	// data := index.Metadata{
	// 	"target":  "car",
	// 	"color":   "red",
	// 	"license": "120",
	// 	"address": "http://iot.qingcloud.com/test.jpg",
	// }
	// data := index.Metadata{
	// 	"target": "left",
	// 	"cause":  "fail",
	// }
	// for {
	// 	reply, err := m.PubEventSync(context.Background(), "FailureReport", data)
	// 	assert.Nil(t, err)
	// 	fmt.Println(reply)
	// 	time.Sleep(1 * time.Second)
	// }
	go func() {
		var i int64 = 1539362482000
		for {
			data := index.Metadata{
				"CO2Concentration": RandInt64(1, 100),
				"humidity":         RandInt64(1, 100),
			}
			//tm := (time.Now().Unix() - int64(i*60*60)) * 1000
			ch, err := m.PubPropertyAsync(data)
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
			time.Sleep(10000 * time.Millisecond)
		}
		//os.Exit(0)
	}()
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

// ----------------------------------------------------------------

/*
	发送消息的几种方式: async sync event
	async: 异步，借助 cache 的handler，t 为当前时间
	sync: 同步，t 为入参
	event: 整个 meta 塞进结构体
*/

func TestPubPropertyAsync(t *testing.T) {
	options := &index.Options{
		Token: "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY3IiOiIxIiwiYXVkIjoiaWFtIiwiYXpwIjoiaWFtIiwiY3VpZCI6ImlhbXItaDJjMjl2OXIiLCJlaXNrIjoiVjRSd3NoNjRpcXJhSTVJTHlnZ2xHZFhnV3E1S1JGWWxFYnRwakkxZk9Raz0iLCJleHAiOjE2MDQ1NDc0ODYsImlhdCI6MTU3MzAxMTQ4NiwiaXNzIjoic3RzIiwianRpIjoiVWpJNFdQQW9wNWdQNldPdHJIUTRlQiIsIm5iZiI6MCwib3JnaSI6ImlvdGQtMjk5MDY2MDktNzNiYS00NzBkLWE2ZmQtMGUxYzE3MTkwZmQwIiwib3d1ciI6InVzci1rTFZWQkRxZCIsInByZWYiOiJxcm46cWluZ2Nsb3VkOmlhbToiLCJydHlwIjoicm9sZSIsInN1YiI6InN0cyIsInRoaWQiOiJpb3R0LWZUeXRjS1BWTlEiLCJ0eXAiOiJJRCJ9.C1oCwaviLAsmb42mDXX4mXw2h0ccXYV8Kd5mAGkCpxpOFM7Rd7lOL2kGMJpvv_I5caOTlSNiFwMe2L2eXiA_dsZPBEW08dmzghLZXpVABFG7KJOrxT5t6WBYzVCOezt4CynSXheIs0NjSMZ5VBTdiEjj8GIi5iAIWUaYrEeFOlj3IZPp7ddr82rkog9OIDnHDvyXDK2MruKAb7xZ2QZFa0Wg1GKixFUhfT0iU37pQZbsGAduj-kB9z4o_ZwtP8gFko6AkW8WuBzzXhs35cQty2vXJ3ohxKnXtoiwChNfIQmNr8Cc7VJmQTmrQPrgmK3uMnxi02SQXsF2vd0HmpA_7A",
		//DeviceId: "iotd-78c675f2-3495-4734-af5b-bb31f83f764c",
		Server: "tcp://127.0.0.1:1883", // tcp://192.168.14.120:1883 tcp://192.168.14.120:8055
		SetProperty: func(meta index.Metadata) (index.Metadata, error) {
			fmt.Println("SetProperty", meta)
			data := make(index.Metadata)
			return data, nil
		},
		ServiceHandle: func(name string, meta index.Metadata) (index.Metadata, error) {
			fmt.Println("ServiceHandle", name, meta)
			data := make(index.Metadata)
			data["test"] = "xxxx"
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

			/*
				data =>	message{
							id
							version
							params: map[data_k]Property{
								data_value
								time
							}
						}
			*/

			ch, err := m.PubPropertyAsync(data)
			// ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			// reply, err := m.PubPropertySync(ctx, data)
			// cancel()
			assert.Nil(t, err)
			select {
			case value := <-ch:
				fmt.Println("expired value:", value)
			}
			// data = index.Metadata{
			// 	"int32":  10,
			// 	"string": "hexing-string",
			// 	"float":  rand.Float32(),
			// 	"double": rand.Float64(),
			// }
			// reply, err = m.PubEventSync(context.Background(), "he-event1", data)
			// assert.Nil(t, err)
			// fmt.Println(reply)
			i = i + 60000
			if i > 1570898482000 {
				i = 1539362482000
			}
			time.Sleep(10000 * time.Millisecond)
		}
		//os.Exit(0)
	}()

	select {}
	//name := "test"
	//reply = m.PubEvent(ctx, name, data)
	//t.Log(reply)
}

func TestPubPropertySync(t *testing.T) {
	options := &index.Options{
		Token: "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY3IiOiIxIiwiYXVkIjoiaWFtIiwiYXpwIjoiaWFtIiwiY3VpZCI6ImlhbXItbDZjdGJoZjAiLCJlaXNrIjoiR2VuazUxbm5BLXZyOUJaSnJQQ1gwNnNPSnBabElFZmw4eGlkVEFNbWRjQT0iLCJleHAiOjE2MDU5Njg2MTQsImlhdCI6MTU3NDQzMjYxNCwiaXNzIjoic3RzIiwianRpIjoiVWpJNFdQQW9wNWdQNldPdHJIUTc5USIsIm5iZiI6MCwib3JnaSI6ImlvdGQtZDhjYmEzOTItYWU0NC00MGRmLTk2YzgtNmQ3MWMzMmI4NjZlIiwib3d1ciI6InVzci16eUt5UFNmRyIsInByZWYiOiJxcm46cWluZ2Nsb3VkOmlhbToiLCJydHlwIjoicm9sZSIsInN1YiI6InN0cyIsInRoaWQiOiJpb3R0LVdDcnQ5bk1hUFMiLCJ0eXAiOiJJRCJ9.V0hqewKk6cwwlWzUpBY1HFpMcEvElurmKHh_HtAD816oVsEvl58kK4zpfs1jslASfBLw11OHBE-BD1Zp9FfGicRgTulQ2OUI4t9UiDbmnxGGKODknuP-0lEAb30n6JqLWWZh-rlZlN0tQVixelMC45ftf4LR0OmRH1T250RWO1MNNqqNgral9juTZ8mI9qcvX0yN3Ro7hM_JndeFWc4j9uj_QLus-Sv0mhleMh4i_5uoji7p8XReykwC82Lm2o61EGZZ3T7RCW9GCrSFngIsXnFUxk9mGqUiyW4aqKNkvpcCg-lm3t4fuszc6YW9_YzU53uic14ERRswREf3Wj3vJg",
		//DeviceId: "iotd-78c675f2-3495-4734-af5b-bb31f83f764c",
		Server: "tcp://127.0.0.1:1883", //192.168.14.120:1883
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

			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			reply, err := m.PubPropertySync(ctx, data)
			cancel()
			assert.Nil(t, err)

			fmt.Println("PubPropertySync reply", reply)

			i = i + 60000
			if i > 1570898482000 {
				i = 1539362482000
			}
			time.Sleep(10000 * time.Millisecond)
		}
	}()
	select {}
}

func TestPubEventSync(t *testing.T) {
	options := &index.Options{
		Token: "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY3IiOiIxIiwiYXVkIjoiaWFtIiwiYXpwIjoiaWFtIiwiY3VpZCI6ImlhbXItbDZjdGJoZjAiLCJlaXNrIjoiR2VuazUxbm5BLXZyOUJaSnJQQ1gwNnNPSnBabElFZmw4eGlkVEFNbWRjQT0iLCJleHAiOjE2MDU5Njg2MTQsImlhdCI6MTU3NDQzMjYxNCwiaXNzIjoic3RzIiwianRpIjoiVWpJNFdQQW9wNWdQNldPdHJIUTc5USIsIm5iZiI6MCwib3JnaSI6ImlvdGQtZDhjYmEzOTItYWU0NC00MGRmLTk2YzgtNmQ3MWMzMmI4NjZlIiwib3d1ciI6InVzci16eUt5UFNmRyIsInByZWYiOiJxcm46cWluZ2Nsb3VkOmlhbToiLCJydHlwIjoicm9sZSIsInN1YiI6InN0cyIsInRoaWQiOiJpb3R0LVdDcnQ5bk1hUFMiLCJ0eXAiOiJJRCJ9.V0hqewKk6cwwlWzUpBY1HFpMcEvElurmKHh_HtAD816oVsEvl58kK4zpfs1jslASfBLw11OHBE-BD1Zp9FfGicRgTulQ2OUI4t9UiDbmnxGGKODknuP-0lEAb30n6JqLWWZh-rlZlN0tQVixelMC45ftf4LR0OmRH1T250RWO1MNNqqNgral9juTZ8mI9qcvX0yN3Ro7hM_JndeFWc4j9uj_QLus-Sv0mhleMh4i_5uoji7p8XReykwC82Lm2o61EGZZ3T7RCW9GCrSFngIsXnFUxk9mGqUiyW4aqKNkvpcCg-lm3t4fuszc6YW9_YzU53uic14ERRswREf3Wj3vJg",
		//DeviceId: "iotd-78c675f2-3495-4734-af5b-bb31f83f764c",
		Server: "tcp://127.0.0.1:1883", // 192.168.14.120:1883
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
				"int32":  10,
				"string": "hexing-string",
				"float":  rand.Float32(),
				"double": rand.Float64(),
			}
			reply, err := m.PubEventSync(context.Background(), "he-event1", data)
			assert.Nil(t, err)
			fmt.Println(reply)
			i = i + 60000
			if i > 1570898482000 {
				i = 1539362482000
			}
			time.Sleep(10000 * time.Millisecond)
		}
	}()
	select {}
}

func TestPubEventAsync(t *testing.T) {
	options := &index.Options{
		Token:  "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY3IiOiIxIiwiYXVkIjoiaWFtIiwiYXpwIjoiaWFtIiwiY3VpZCI6ImlhbXItaDJjMjl2OXIiLCJlaXNrIjoiVjRSd3NoNjRpcXJhSTVJTHlnZ2xHZFhnV3E1S1JGWWxFYnRwakkxZk9Raz0iLCJleHAiOjE2MDQ1NDc0ODYsImlhdCI6MTU3MzAxMTQ4NiwiaXNzIjoic3RzIiwianRpIjoiVWpJNFdQQW9wNWdQNldPdHJIUTRlQiIsIm5iZiI6MCwib3JnaSI6ImlvdGQtMjk5MDY2MDktNzNiYS00NzBkLWE2ZmQtMGUxYzE3MTkwZmQwIiwib3d1ciI6InVzci1rTFZWQkRxZCIsInByZWYiOiJxcm46cWluZ2Nsb3VkOmlhbToiLCJydHlwIjoicm9sZSIsInN1YiI6InN0cyIsInRoaWQiOiJpb3R0LWZUeXRjS1BWTlEiLCJ0eXAiOiJJRCJ9.C1oCwaviLAsmb42mDXX4mXw2h0ccXYV8Kd5mAGkCpxpOFM7Rd7lOL2kGMJpvv_I5caOTlSNiFwMe2L2eXiA_dsZPBEW08dmzghLZXpVABFG7KJOrxT5t6WBYzVCOezt4CynSXheIs0NjSMZ5VBTdiEjj8GIi5iAIWUaYrEeFOlj3IZPp7ddr82rkog9OIDnHDvyXDK2MruKAb7xZ2QZFa0Wg1GKixFUhfT0iU37pQZbsGAduj-kB9z4o_ZwtP8gFko6AkW8WuBzzXhs35cQty2vXJ3ohxKnXtoiwChNfIQmNr8Cc7VJmQTmrQPrgmK3uMnxi02SQXsF2vd0HmpA_7A",
		Server: "tcp://127.0.0.1:1883", // tcp://192.168.14.120:1883 tcp://192.168.14.120:8055
		SetProperty: func(meta index.Metadata) (index.Metadata, error) {
			fmt.Println("SetProperty", meta)
			data := make(index.Metadata)
			return data, nil
		},
		ServiceHandle: func(name string, meta index.Metadata) (index.Metadata, error) {
			fmt.Println("ServiceHandle", name, meta)
			data := make(index.Metadata)
			data["test"] = "xxxx"
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

			ch, err := m.PubEventAsync("he-event1", data)
			assert.Nil(t, err)
			select {
			case value := <-ch:
				fmt.Println("expired value:", value)
			}

			i = i + 60000
			if i > 1570898482000 {
				i = 1539362482000
			}
			time.Sleep(10000 * time.Millisecond)
		}
	}()

	select {}
}

// sub
func TestSubDeviceControlAsync(t *testing.T) {
	options := &index.Options{
		Token:  "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY3IiOiIxIiwiYXVkIjoiaWFtIiwiYXpwIjoiaWFtIiwiY3VpZCI6ImlhbXItaDJjMjl2OXIiLCJlaXNrIjoiVjRSd3NoNjRpcXJhSTVJTHlnZ2xHZFhnV3E1S1JGWWxFYnRwakkxZk9Raz0iLCJleHAiOjE2MDQ1NDc0ODYsImlhdCI6MTU3MzAxMTQ4NiwiaXNzIjoic3RzIiwianRpIjoiVWpJNFdQQW9wNWdQNldPdHJIUTRlQiIsIm5iZiI6MCwib3JnaSI6ImlvdGQtMjk5MDY2MDktNzNiYS00NzBkLWE2ZmQtMGUxYzE3MTkwZmQwIiwib3d1ciI6InVzci1rTFZWQkRxZCIsInByZWYiOiJxcm46cWluZ2Nsb3VkOmlhbToiLCJydHlwIjoicm9sZSIsInN1YiI6InN0cyIsInRoaWQiOiJpb3R0LWZUeXRjS1BWTlEiLCJ0eXAiOiJJRCJ9.C1oCwaviLAsmb42mDXX4mXw2h0ccXYV8Kd5mAGkCpxpOFM7Rd7lOL2kGMJpvv_I5caOTlSNiFwMe2L2eXiA_dsZPBEW08dmzghLZXpVABFG7KJOrxT5t6WBYzVCOezt4CynSXheIs0NjSMZ5VBTdiEjj8GIi5iAIWUaYrEeFOlj3IZPp7ddr82rkog9OIDnHDvyXDK2MruKAb7xZ2QZFa0Wg1GKixFUhfT0iU37pQZbsGAduj-kB9z4o_ZwtP8gFko6AkW8WuBzzXhs35cQty2vXJ3ohxKnXtoiwChNfIQmNr8Cc7VJmQTmrQPrgmK3uMnxi02SQXsF2vd0HmpA_7A",
		Server: "tcp://127.0.0.1:1883", // tcp://192.168.14.120:1883 tcp://192.168.14.120:8055
		SetProperty: func(meta index.Metadata) (index.Metadata, error) {
			fmt.Println("SetProperty", meta)
			data := make(index.Metadata)
			return data, nil
		},
		ServiceHandle: func(name string, meta index.Metadata) (index.Metadata, error) {
			fmt.Println("ServiceHandle", name, meta)
			data := make(index.Metadata)
			data["test"] = "xxxx"
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
			reply, err := m.SubDeviceControlSync("he-event1")
			assert.Nil(t, err)

			fmt.Println("=====reply:", reply)

			i = i + 60000
			if i > 1570898482000 {
				i = 1539362482000
			}
			time.Sleep(10000 * time.Millisecond)
		}
	}()

	select {}
}
