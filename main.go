package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"git.internal.yunify.com/iot-sdk/device-sdk-go/index"
	"git.internal.yunify.com/iot-sdk/device-sdk-go/mqtt"
)

func main() {
	// 设备上线
	// deviceConnect()

	// 上报设备属性
	// pubDeviceProperty()

	// 上报设备事件
	pubDeviceEvent()

	// 设备控制
	// deviceControl()
}

// 设备连接
func deviceConnect() {
	options := &index.Options{
		Token:  "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY3IiOiIxIiwiYXVkIjoiaWFtIiwiYXpwIjoiaWFtIiwiY3VpZCI6ImlhbXItdW9nbjh3N3IiLCJlaXNrIjoiSi1hSHdySDgwTHVXWEI1bGJMQ2E3cm1uTFVfZ1RlbWhJd0FJdVN6T29qZz0iLCJleHAiOjE2MTI5MTk3MjAsImlhdCI6MTU4MTM4MzcyMCwiaXNzIjoic3RzIiwianRpIjoiVWpJNFdQQW9wNWdQNldPdHJIUU5rUyIsIm5iZiI6MCwib3JnaSI6ImlvdGQtYjViODJiYzMtZWEwYS00NDRjLTliYjMtNDkzOWE0YTgzMmRiIiwib3d1ciI6InVzci1rZUF5dG16MSIsInByZWYiOiJxcm46cWluZ2Nsb3VkOmlhbToiLCJydHlwIjoicm9sZSIsInN1YiI6InN0cyIsInRoaWQiOiJpb3R0LXlWQXd4OXJiOGoiLCJ0eXAiOiJJRCJ9.Hq6zIwQCpBV897fsVl-WHxBtgtH8xe8umcp5QIQ3p1lSHrYUV_ofrbJ5oZKasUKwYqxlhhcxjX2f3U9OCCvOj8yGjIyK8vrf8vJBbNwW48fkCiVnFOpoKui8k9Fg13qNl0AUD8TmOWAukn3uQTI7gKW6fhwmWkdZD8cLOraEBvkGkrL19Nlw-JuU-MWXeNB2p1F5CahUAvDD78zUHkPJTZ-X3v9d73YyUlSV2CrJvAJGpae6sHCXk4iS8KyPQ1GNjPmjD9qBdbzr5cdIA3LjIkuppaWb0i8vymhvLaqcfD5EnEfu8aKNNLGBedEI3c8BlXOLSgp5_BldOJuP2GnGfQ",
		Server: "tcp://192.168.14.120:8055", // 127.0.0.1:1883
	}

	m, err := mqtt.NewMqtt(options)
	if err != nil {
		panic(err)
	}

	// 连接
	err = m.Connect()
	if err != nil {
		panic(err)
	}
}

func pubDeviceProperty() {
	options := &index.Options{
		Token:  "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY3IiOiIxIiwiYXVkIjoiaWFtIiwiYXpwIjoiaWFtIiwiY3VpZCI6ImlhbXItdW9nbjh3N3IiLCJlaXNrIjoiSi1hSHdySDgwTHVXWEI1bGJMQ2E3cm1uTFVfZ1RlbWhJd0FJdVN6T29qZz0iLCJleHAiOjE2MTI5MTk3MjAsImlhdCI6MTU4MTM4MzcyMCwiaXNzIjoic3RzIiwianRpIjoiVWpJNFdQQW9wNWdQNldPdHJIUU5rUyIsIm5iZiI6MCwib3JnaSI6ImlvdGQtYjViODJiYzMtZWEwYS00NDRjLTliYjMtNDkzOWE0YTgzMmRiIiwib3d1ciI6InVzci1rZUF5dG16MSIsInByZWYiOiJxcm46cWluZ2Nsb3VkOmlhbToiLCJydHlwIjoicm9sZSIsInN1YiI6InN0cyIsInRoaWQiOiJpb3R0LXlWQXd4OXJiOGoiLCJ0eXAiOiJJRCJ9.Hq6zIwQCpBV897fsVl-WHxBtgtH8xe8umcp5QIQ3p1lSHrYUV_ofrbJ5oZKasUKwYqxlhhcxjX2f3U9OCCvOj8yGjIyK8vrf8vJBbNwW48fkCiVnFOpoKui8k9Fg13qNl0AUD8TmOWAukn3uQTI7gKW6fhwmWkdZD8cLOraEBvkGkrL19Nlw-JuU-MWXeNB2p1F5CahUAvDD78zUHkPJTZ-X3v9d73YyUlSV2CrJvAJGpae6sHCXk4iS8KyPQ1GNjPmjD9qBdbzr5cdIA3LjIkuppaWb0i8vymhvLaqcfD5EnEfu8aKNNLGBedEI3c8BlXOLSgp5_BldOJuP2GnGfQ",
		Server: "tcp://127.0.0.1:1883", // 127.0.0.1:1883
	}

	m, err := mqtt.NewMqtt(options)
	if err != nil {
		panic(err)
	}

	// 连接
	err = m.Connect()
	if err != nil {
		panic(err)
	}

	data := index.Metadata{
		"CO2Concentration": mqtt.RandInt64(1, 100),
		"humidity":         mqtt.RandInt64(1, 100),
	}

	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
	reply, err := m.PubProperty(ctx, data)
	if err != nil {
		panic(err)
	}

	fmt.Println("PubPropertySync reply", reply)
}

func pubDeviceEvent() {
	options := &index.Options{
		Token:  "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY3IiOiIxIiwiYXVkIjoiaWFtIiwiYXpwIjoiaWFtIiwiY3VpZCI6ImlhbXItdW9nbjh3N3IiLCJlaXNrIjoiSi1hSHdySDgwTHVXWEI1bGJMQ2E3cm1uTFVfZ1RlbWhJd0FJdVN6T29qZz0iLCJleHAiOjE2MTI5MTk3MjAsImlhdCI6MTU4MTM4MzcyMCwiaXNzIjoic3RzIiwianRpIjoiVWpJNFdQQW9wNWdQNldPdHJIUU5rUyIsIm5iZiI6MCwib3JnaSI6ImlvdGQtYjViODJiYzMtZWEwYS00NDRjLTliYjMtNDkzOWE0YTgzMmRiIiwib3d1ciI6InVzci1rZUF5dG16MSIsInByZWYiOiJxcm46cWluZ2Nsb3VkOmlhbToiLCJydHlwIjoicm9sZSIsInN1YiI6InN0cyIsInRoaWQiOiJpb3R0LXlWQXd4OXJiOGoiLCJ0eXAiOiJJRCJ9.Hq6zIwQCpBV897fsVl-WHxBtgtH8xe8umcp5QIQ3p1lSHrYUV_ofrbJ5oZKasUKwYqxlhhcxjX2f3U9OCCvOj8yGjIyK8vrf8vJBbNwW48fkCiVnFOpoKui8k9Fg13qNl0AUD8TmOWAukn3uQTI7gKW6fhwmWkdZD8cLOraEBvkGkrL19Nlw-JuU-MWXeNB2p1F5CahUAvDD78zUHkPJTZ-X3v9d73YyUlSV2CrJvAJGpae6sHCXk4iS8KyPQ1GNjPmjD9qBdbzr5cdIA3LjIkuppaWb0i8vymhvLaqcfD5EnEfu8aKNNLGBedEI3c8BlXOLSgp5_BldOJuP2GnGfQ",
		Server: "tcp://127.0.0.1:1883", // 127.0.0.1:1883  192.168.14.120:8055
	}

	m, err := mqtt.NewMqtt(options)
	if err != nil {
		panic(err)
	}

	// 连接
	err = m.Connect()
	if err != nil {
		panic(err)
	}

	data := index.Metadata{
		"int32":  10,
		"string": "mj-string",
		"float":  rand.Float32(),
		"double": rand.Float64(),
	}
	reply, err := m.PubEvent(context.Background(), "mj-event1", data)
	if err != nil {
		panic(err)
	}
	fmt.Printf("PubEvent reply:%+v\n", reply)

}

func deviceControl() {
	options := &index.Options{
		Token:  "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY3IiOiIxIiwiYXVkIjoiaWFtIiwiYXpwIjoiaWFtIiwiY3VpZCI6ImlhbXItdW9nbjh3N3IiLCJlaXNrIjoiSi1hSHdySDgwTHVXWEI1bGJMQ2E3cm1uTFVfZ1RlbWhJd0FJdVN6T29qZz0iLCJleHAiOjE2MTI5MTk3MjAsImlhdCI6MTU4MTM4MzcyMCwiaXNzIjoic3RzIiwianRpIjoiVWpJNFdQQW9wNWdQNldPdHJIUU5rUyIsIm5iZiI6MCwib3JnaSI6ImlvdGQtYjViODJiYzMtZWEwYS00NDRjLTliYjMtNDkzOWE0YTgzMmRiIiwib3d1ciI6InVzci1rZUF5dG16MSIsInByZWYiOiJxcm46cWluZ2Nsb3VkOmlhbToiLCJydHlwIjoicm9sZSIsInN1YiI6InN0cyIsInRoaWQiOiJpb3R0LXlWQXd4OXJiOGoiLCJ0eXAiOiJJRCJ9.Hq6zIwQCpBV897fsVl-WHxBtgtH8xe8umcp5QIQ3p1lSHrYUV_ofrbJ5oZKasUKwYqxlhhcxjX2f3U9OCCvOj8yGjIyK8vrf8vJBbNwW48fkCiVnFOpoKui8k9Fg13qNl0AUD8TmOWAukn3uQTI7gKW6fhwmWkdZD8cLOraEBvkGkrL19Nlw-JuU-MWXeNB2p1F5CahUAvDD78zUHkPJTZ-X3v9d73YyUlSV2CrJvAJGpae6sHCXk4iS8KyPQ1GNjPmjD9qBdbzr5cdIA3LjIkuppaWb0i8vymhvLaqcfD5EnEfu8aKNNLGBedEI3c8BlXOLSgp5_BldOJuP2GnGfQ",
		Server: "tcp://192.168.14.120:8055", // 127.0.0.1:1883
	}

	m, err := mqtt.NewMqtt(options)
	if err != nil {
		panic(err)
	}

	// 连接
	err = m.Connect()
	if err != nil {
		panic(err)
	}

	go m.SubDeviceControl()

	time.Sleep(15 * time.Second)

	m.UnSubDeviceControl()

	time.Sleep(3 * time.Second)
}
