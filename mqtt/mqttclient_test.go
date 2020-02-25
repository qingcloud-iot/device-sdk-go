package mqtt

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"git.internal.yunify.com/iot-sdk/device-sdk-go/index"
	cache "github.com/muesli/cache2go"
	"github.com/stretchr/testify/assert"
)

func TestNewMqtt(t *testing.T) {
	options := &index.Options{
		Token: "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY3IiOiIxIiwiYXVkIjoiaWFtIiwiYXpwIjoiaWFtIiwiY3VpZCI6ImlhbXItdW9nbjh3N3IiLCJlaXNrIjoiSi1hSHdySDgwTHVXWEI1bGJMQ2E3cm1uTFVfZ1RlbWhJd0FJdVN6T29qZz0iLCJleHAiOjE2MTI5MTk3MjAsImlhdCI6MTU4MTM4MzcyMCwiaXNzIjoic3RzIiwianRpIjoiVWpJNFdQQW9wNWdQNldPdHJIUU5rUyIsIm5iZiI6MCwib3JnaSI6ImlvdGQtYjViODJiYzMtZWEwYS00NDRjLTliYjMtNDkzOWE0YTgzMmRiIiwib3d1ciI6InVzci1rZUF5dG16MSIsInByZWYiOiJxcm46cWluZ2Nsb3VkOmlhbToiLCJydHlwIjoicm9sZSIsInN1YiI6InN0cyIsInRoaWQiOiJpb3R0LXlWQXd4OXJiOGoiLCJ0eXAiOiJJRCJ9.Hq6zIwQCpBV897fsVl-WHxBtgtH8xe8umcp5QIQ3p1lSHrYUV_ofrbJ5oZKasUKwYqxlhhcxjX2f3U9OCCvOj8yGjIyK8vrf8vJBbNwW48fkCiVnFOpoKui8k9Fg13qNl0AUD8TmOWAukn3uQTI7gKW6fhwmWkdZD8cLOraEBvkGkrL19Nlw-JuU-MWXeNB2p1F5CahUAvDD78zUHkPJTZ-X3v9d73YyUlSV2CrJvAJGpae6sHCXk4iS8KyPQ1GNjPmjD9qBdbzr5cdIA3LjIkuppaWb0i8vymhvLaqcfD5EnEfu8aKNNLGBedEI3c8BlXOLSgp5_BldOJuP2GnGfQ",
		//entityID: "iotd-78c675f2-3495-4734-af5b-bb31f83f764c",
		Server: "tcp://192.168.14.120:8055", // 192.168.14.120:8055
	}
	m, err := NewMqtt(options)
	assert.Nil(t, err)
	err = m.Connect()
	assert.Nil(t, err)
	time.Sleep(5 * time.Second)

	// go func() {
	// 	var i int64 = 1539362482000
	// 	for {
	// 		data := index.Metadata{
	// 			"CO2Concentration": RandInt64(1, 100),
	// 			"humidity":         RandInt64(1, 100),
	// 		}
	// 		ch, err := m.PubPropertyAsync(data)
	// 		assert.Nil(t, err)

	// 		select {
	// 		case value := <-ch:
	// 			fmt.Println(value)
	// 		}

	// 		i = i + 60000
	// 		if i > 1570898482000 {
	// 			i = 1539362482000
	// 		}
	// 		time.Sleep(10000 * time.Millisecond)
	// 	}
	// }()
	select {}
}

func TestTask(t *testing.T) {
	cache := cache.Cache("xxxx")
	cache.Add("a", 3*time.Second, "xxx")
	cache.RemoveAboutToDeleteItemCallback()
}

// ----------------------------------------------------------------

/*
	发送消息的几种方式: async sync，一般用 sync
	async: 异步，借助 cache 的handler，t 为当前时间
	sync: 同步，t 为入参
	event: 整个 meta 塞进结构体
*/

func TestPubPropertyAsync(t *testing.T) {
	options := &index.Options{
		Token: "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY3IiOiIxIiwiYXVkIjoiaWFtIiwiYXpwIjoiaWFtIiwiY3VpZCI6ImlhbXItaDJjMjl2OXIiLCJlaXNrIjoiVjRSd3NoNjRpcXJhSTVJTHlnZ2xHZFhnV3E1S1JGWWxFYnRwakkxZk9Raz0iLCJleHAiOjE2MDQ1NDc0ODYsImlhdCI6MTU3MzAxMTQ4NiwiaXNzIjoic3RzIiwianRpIjoiVWpJNFdQQW9wNWdQNldPdHJIUTRlQiIsIm5iZiI6MCwib3JnaSI6ImlvdGQtMjk5MDY2MDktNzNiYS00NzBkLWE2ZmQtMGUxYzE3MTkwZmQwIiwib3d1ciI6InVzci1rTFZWQkRxZCIsInByZWYiOiJxcm46cWluZ2Nsb3VkOmlhbToiLCJydHlwIjoicm9sZSIsInN1YiI6InN0cyIsInRoaWQiOiJpb3R0LWZUeXRjS1BWTlEiLCJ0eXAiOiJJRCJ9.C1oCwaviLAsmb42mDXX4mXw2h0ccXYV8Kd5mAGkCpxpOFM7Rd7lOL2kGMJpvv_I5caOTlSNiFwMe2L2eXiA_dsZPBEW08dmzghLZXpVABFG7KJOrxT5t6WBYzVCOezt4CynSXheIs0NjSMZ5VBTdiEjj8GIi5iAIWUaYrEeFOlj3IZPp7ddr82rkog9OIDnHDvyXDK2MruKAb7xZ2QZFa0Wg1GKixFUhfT0iU37pQZbsGAduj-kB9z4o_ZwtP8gFko6AkW8WuBzzXhs35cQty2vXJ3ohxKnXtoiwChNfIQmNr8Cc7VJmQTmrQPrgmK3uMnxi02SQXsF2vd0HmpA_7A",
		//entityID: "iotd-78c675f2-3495-4734-af5b-bb31f83f764c",
		Server: "tcp://127.0.0.1:1883", // tcp://192.168.14.120:1883 tcp://192.168.14.120:8055
	}
	m, err := NewMqtt(options)
	assert.Nil(t, err)
	assert.Nil(t, m)
	time.Sleep(5 * time.Second)

	go func() {
		var i int64 = 1539362482000
		for {
			data := index.PropertyKV{
				"CO2Concentration": RandInt64(1, 100),
				"humidity":         RandInt64(1, 100),
			}

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

func TestPubPropertySync(t *testing.T) {
	options := &index.Options{
		Token:  "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY3IiOiIxIiwiYXVkIjoiaWFtIiwiYXpwIjoiaWFtIiwiY3VpZCI6ImlhbXItbDZjdGJoZjAiLCJlaXNrIjoiR2VuazUxbm5BLXZyOUJaSnJQQ1gwNnNPSnBabElFZmw4eGlkVEFNbWRjQT0iLCJleHAiOjE2MDU5Njg2MTQsImlhdCI6MTU3NDQzMjYxNCwiaXNzIjoic3RzIiwianRpIjoiVWpJNFdQQW9wNWdQNldPdHJIUTc5USIsIm5iZiI6MCwib3JnaSI6ImlvdGQtZDhjYmEzOTItYWU0NC00MGRmLTk2YzgtNmQ3MWMzMmI4NjZlIiwib3d1ciI6InVzci16eUt5UFNmRyIsInByZWYiOiJxcm46cWluZ2Nsb3VkOmlhbToiLCJydHlwIjoicm9sZSIsInN1YiI6InN0cyIsInRoaWQiOiJpb3R0LVdDcnQ5bk1hUFMiLCJ0eXAiOiJJRCJ9.V0hqewKk6cwwlWzUpBY1HFpMcEvElurmKHh_HtAD816oVsEvl58kK4zpfs1jslASfBLw11OHBE-BD1Zp9FfGicRgTulQ2OUI4t9UiDbmnxGGKODknuP-0lEAb30n6JqLWWZh-rlZlN0tQVixelMC45ftf4LR0OmRH1T250RWO1MNNqqNgral9juTZ8mI9qcvX0yN3Ro7hM_JndeFWc4j9uj_QLus-Sv0mhleMh4i_5uoji7p8XReykwC82Lm2o61EGZZ3T7RCW9GCrSFngIsXnFUxk9mGqUiyW4aqKNkvpcCg-lm3t4fuszc6YW9_YzU53uic14ERRswREf3Wj3vJg",
		Server: "tcp://127.0.0.1:1883", //192.168.14.120:1883
	}

	m, err := NewMqtt(options)
	assert.Nil(t, err)
	assert.Nil(t, m)
	time.Sleep(5 * time.Second)

	go func() {
		var i int64 = 1539362482000
		for {
			data := index.PropertyKV{
				"CO2Concentration": RandInt64(1, 100),
				"humidity":         RandInt64(1, 100),
			}

			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			reply, err := m.PubProperty(ctx, data)
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
		Token:  "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY3IiOiIxIiwiYXVkIjoiaWFtIiwiYXpwIjoiaWFtIiwiY3VpZCI6ImlhbXItbDZjdGJoZjAiLCJlaXNrIjoiR2VuazUxbm5BLXZyOUJaSnJQQ1gwNnNPSnBabElFZmw4eGlkVEFNbWRjQT0iLCJleHAiOjE2MDU5Njg2MTQsImlhdCI6MTU3NDQzMjYxNCwiaXNzIjoic3RzIiwianRpIjoiVWpJNFdQQW9wNWdQNldPdHJIUTc5USIsIm5iZiI6MCwib3JnaSI6ImlvdGQtZDhjYmEzOTItYWU0NC00MGRmLTk2YzgtNmQ3MWMzMmI4NjZlIiwib3d1ciI6InVzci16eUt5UFNmRyIsInByZWYiOiJxcm46cWluZ2Nsb3VkOmlhbToiLCJydHlwIjoicm9sZSIsInN1YiI6InN0cyIsInRoaWQiOiJpb3R0LVdDcnQ5bk1hUFMiLCJ0eXAiOiJJRCJ9.V0hqewKk6cwwlWzUpBY1HFpMcEvElurmKHh_HtAD816oVsEvl58kK4zpfs1jslASfBLw11OHBE-BD1Zp9FfGicRgTulQ2OUI4t9UiDbmnxGGKODknuP-0lEAb30n6JqLWWZh-rlZlN0tQVixelMC45ftf4LR0OmRH1T250RWO1MNNqqNgral9juTZ8mI9qcvX0yN3Ro7hM_JndeFWc4j9uj_QLus-Sv0mhleMh4i_5uoji7p8XReykwC82Lm2o61EGZZ3T7RCW9GCrSFngIsXnFUxk9mGqUiyW4aqKNkvpcCg-lm3t4fuszc6YW9_YzU53uic14ERRswREf3Wj3vJg",
		Server: "tcp://127.0.0.1:1883", // 192.168.14.120:1883
	}
	m, err := NewMqtt(options)
	assert.Nil(t, err)
	assert.Nil(t, m)
	time.Sleep(5 * time.Second)

	go func() {
		var i int64 = 1539362482000
		for {
			data := index.PropertyKV{
				"int32":  10,
				"string": "hexing-string",
				"float":  rand.Float32(),
				"double": rand.Float64(),
			}
			reply, err := m.PubEvent(context.Background(), "he-event1", data)
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
	}
	m, err := NewMqtt(options)
	assert.Nil(t, err)
	assert.Nil(t, m)
	time.Sleep(5 * time.Second)

	go func() {
		var i int64 = 1539362482000
		for {
			data := index.PropertyKV{
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
/*
	运行 ehub
	运行 mqttbox
	在 mqttbox 中向 topic 发布消息，订阅 topic_reply 消息
	{
    "id": "123",
    "version": "1.0",
    "params": {
        "label":"on",
        "image":"23.6"
    	}
	}

	获得：
	{
    "code":200,
    "id":"123",
    "data":{
        "image":"23.6",
        "label":"on"
    	}
	}
*/
func TestSubDeviceControlSync(t *testing.T) {
	options := &index.Options{
		Token:     "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY3IiOiIxIiwiYXVkIjoiaWFtIiwiYXpwIjoiaWFtIiwiY3VpZCI6ImlhbXItaDJjMjl2OXIiLCJlaXNrIjoiVjRSd3NoNjRpcXJhSTVJTHlnZ2xHZFhnV3E1S1JGWWxFYnRwakkxZk9Raz0iLCJleHAiOjE2MDQ1NDc0ODYsImlhdCI6MTU3MzAxMTQ4NiwiaXNzIjoic3RzIiwianRpIjoiVWpJNFdQQW9wNWdQNldPdHJIUTRlQiIsIm5iZiI6MCwib3JnaSI6ImlvdGQtMjk5MDY2MDktNzNiYS00NzBkLWE2ZmQtMGUxYzE3MTkwZmQwIiwib3d1ciI6InVzci1rTFZWQkRxZCIsInByZWYiOiJxcm46cWluZ2Nsb3VkOmlhbToiLCJydHlwIjoicm9sZSIsInN1YiI6InN0cyIsInRoaWQiOiJpb3R0LWZUeXRjS1BWTlEiLCJ0eXAiOiJJRCJ9.C1oCwaviLAsmb42mDXX4mXw2h0ccXYV8Kd5mAGkCpxpOFM7Rd7lOL2kGMJpvv_I5caOTlSNiFwMe2L2eXiA_dsZPBEW08dmzghLZXpVABFG7KJOrxT5t6WBYzVCOezt4CynSXheIs0NjSMZ5VBTdiEjj8GIi5iAIWUaYrEeFOlj3IZPp7ddr82rkog9OIDnHDvyXDK2MruKAb7xZ2QZFa0Wg1GKixFUhfT0iU37pQZbsGAduj-kB9z4o_ZwtP8gFko6AkW8WuBzzXhs35cQty2vXJ3ohxKnXtoiwChNfIQmNr8Cc7VJmQTmrQPrgmK3uMnxi02SQXsF2vd0HmpA_7A",
		Server:    "tcp://127.0.0.1:1883", // tcp://192.168.14.120:1883 tcp://192.168.14.120:8055
		Identifer: "start",
	}
	m, err := NewMqtt(options)
	assert.Nil(t, err)

	go m.SubDeviceControl()
	fmt.Println("run")

	time.Sleep(15 * time.Second)

	m.UnSubDeviceControl()

	time.Sleep(3 * time.Second)
}
