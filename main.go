package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
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
	// pubDeviceEvent()

	// 设备控制
	// deviceControl()
}

// 设备连接
func deviceConnect() {
	options := &index.Options{
		Token:  "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY3IiOiIxIiwiYXVkIjoiaWFtIiwiYXpwIjoiaWFtIiwiY3VpZCI6ImlhbXItdm9maW4wYmUiLCJlaXNrIjoiblV4MTJkZDNQWVU1c2RjMlhzcU40Z0I4enNreHVwbTl5R0FjVXFMVDB5az0iLCJleHAiOjE2MTMwMDY3MTUsImlhdCI6MTU4MTQ3MDcxNSwiaXNzIjoic3RzIiwianRpIjoiVWpJNFdQQW9wNWdQNldPdHJIUU82ZCIsIm5iZiI6MCwib3JnaSI6ImlvdGQtNTlmNjg1Y2UtNzBmOS00NDg1LTk5ODUtMjcxZDVkZmI5NDc1Iiwib3d1ciI6InVzci1rZUF5dG16MSIsInByZWYiOiJxcm46cWluZ2Nsb3VkOmlhbToiLCJydHlwIjoicm9sZSIsInN1YiI6InN0cyIsInRoaWQiOiJpb3R0LXlWQXd4OXJiOGoiLCJ0eXAiOiJJRCJ9.M03UZOE_llNCR80LYdmforG5_Bc_QTJN9A2BPLfYX5OZAeawaRoqzOOBIqjORk_HKMLk210ex5DTcQflrUSTNhXiVMilau8a3loi-qY5-13aB45Ra_-qaQpGKcIzCtSsOofNhnOBsshLgvLG0W_ThlY-L5i6FAsTDp9fWKs_hS4VMn1cb8iexi3Oljcy7255J-wWRSaAMcm4KzZNc3kS_HR7NdfGlu9zmjE22rnmlZS60OEvjhqU-SKJBsalHAiFbAWTemHuk5jlB7P2sFiM4JAxIuznq23s0WrNM0oQTRi6xb0bMglGuBmyvPkoh1jMAGklHStprNoxwY_S2aKiUA",
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

	select {}
}

func pubDeviceProperty() {
	options := &index.Options{
		Token:  "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY3IiOiIxIiwiYXVkIjoiaWFtIiwiYXpwIjoiaWFtIiwiY3VpZCI6ImlhbXItdm9maW4wYmUiLCJlaXNrIjoiblV4MTJkZDNQWVU1c2RjMlhzcU40Z0I4enNreHVwbTl5R0FjVXFMVDB5az0iLCJleHAiOjE2MTMwMDY3MTUsImlhdCI6MTU4MTQ3MDcxNSwiaXNzIjoic3RzIiwianRpIjoiVWpJNFdQQW9wNWdQNldPdHJIUU82ZCIsIm5iZiI6MCwib3JnaSI6ImlvdGQtNTlmNjg1Y2UtNzBmOS00NDg1LTk5ODUtMjcxZDVkZmI5NDc1Iiwib3d1ciI6InVzci1rZUF5dG16MSIsInByZWYiOiJxcm46cWluZ2Nsb3VkOmlhbToiLCJydHlwIjoicm9sZSIsInN1YiI6InN0cyIsInRoaWQiOiJpb3R0LXlWQXd4OXJiOGoiLCJ0eXAiOiJJRCJ9.M03UZOE_llNCR80LYdmforG5_Bc_QTJN9A2BPLfYX5OZAeawaRoqzOOBIqjORk_HKMLk210ex5DTcQflrUSTNhXiVMilau8a3loi-qY5-13aB45Ra_-qaQpGKcIzCtSsOofNhnOBsshLgvLG0W_ThlY-L5i6FAsTDp9fWKs_hS4VMn1cb8iexi3Oljcy7255J-wWRSaAMcm4KzZNc3kS_HR7NdfGlu9zmjE22rnmlZS60OEvjhqU-SKJBsalHAiFbAWTemHuk5jlB7P2sFiM4JAxIuznq23s0WrNM0oQTRi6xb0bMglGuBmyvPkoh1jMAGklHStprNoxwY_S2aKiUA",
		Server: "tcp://192.168.14.120:8055", // 127.0.0.1:1883 192.168.14.120:1883
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
		"MaxValue":   float64(22),
		"AlarmState": true,
	}

	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
	reply, err := m.PubProperty(ctx, data)
	if err != nil {
		panic(err)
	}

	fmt.Println("PubPropertySync reply", reply)
	select {}
}

func pubDeviceEvent() {
	options := &index.Options{
		Token:  "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY3IiOiIxIiwiYXVkIjoiaWFtIiwiYXpwIjoiaWFtIiwiY3VpZCI6ImlhbXItdm9maW4wYmUiLCJlaXNrIjoiblV4MTJkZDNQWVU1c2RjMlhzcU40Z0I4enNreHVwbTl5R0FjVXFMVDB5az0iLCJleHAiOjE2MTMwMDY3MTUsImlhdCI6MTU4MTQ3MDcxNSwiaXNzIjoic3RzIiwianRpIjoiVWpJNFdQQW9wNWdQNldPdHJIUU82ZCIsIm5iZiI6MCwib3JnaSI6ImlvdGQtNTlmNjg1Y2UtNzBmOS00NDg1LTk5ODUtMjcxZDVkZmI5NDc1Iiwib3d1ciI6InVzci1rZUF5dG16MSIsInByZWYiOiJxcm46cWluZ2Nsb3VkOmlhbToiLCJydHlwIjoicm9sZSIsInN1YiI6InN0cyIsInRoaWQiOiJpb3R0LXlWQXd4OXJiOGoiLCJ0eXAiOiJJRCJ9.M03UZOE_llNCR80LYdmforG5_Bc_QTJN9A2BPLfYX5OZAeawaRoqzOOBIqjORk_HKMLk210ex5DTcQflrUSTNhXiVMilau8a3loi-qY5-13aB45Ra_-qaQpGKcIzCtSsOofNhnOBsshLgvLG0W_ThlY-L5i6FAsTDp9fWKs_hS4VMn1cb8iexi3Oljcy7255J-wWRSaAMcm4KzZNc3kS_HR7NdfGlu9zmjE22rnmlZS60OEvjhqU-SKJBsalHAiFbAWTemHuk5jlB7P2sFiM4JAxIuznq23s0WrNM0oQTRi6xb0bMglGuBmyvPkoh1jMAGklHStprNoxwY_S2aKiUA",
		Server: "tcp://192.168.14.120:8055", // 127.0.0.1:1883  192.168.14.120:8055
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
		"max": 125.2,
		"min": 10.1,
	}
	reply, err := m.PubEvent(context.Background(), "statistics", data)
	if err != nil {
		panic(err)
	}
	fmt.Printf("PubEvent reply:%+v\n", reply)
	select {}
}

func deviceControl() {
	options := &index.Options{
		Token:     "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY3IiOiIxIiwiYXVkIjoiaWFtIiwiYXpwIjoiaWFtIiwiY3VpZCI6ImlhbXItdm9maW4wYmUiLCJlaXNrIjoiblV4MTJkZDNQWVU1c2RjMlhzcU40Z0I4enNreHVwbTl5R0FjVXFMVDB5az0iLCJleHAiOjE2MTMwMDY3MTUsImlhdCI6MTU4MTQ3MDcxNSwiaXNzIjoic3RzIiwianRpIjoiVWpJNFdQQW9wNWdQNldPdHJIUU82ZCIsIm5iZiI6MCwib3JnaSI6ImlvdGQtNTlmNjg1Y2UtNzBmOS00NDg1LTk5ODUtMjcxZDVkZmI5NDc1Iiwib3d1ciI6InVzci1rZUF5dG16MSIsInByZWYiOiJxcm46cWluZ2Nsb3VkOmlhbToiLCJydHlwIjoicm9sZSIsInN1YiI6InN0cyIsInRoaWQiOiJpb3R0LXlWQXd4OXJiOGoiLCJ0eXAiOiJJRCJ9.M03UZOE_llNCR80LYdmforG5_Bc_QTJN9A2BPLfYX5OZAeawaRoqzOOBIqjORk_HKMLk210ex5DTcQflrUSTNhXiVMilau8a3loi-qY5-13aB45Ra_-qaQpGKcIzCtSsOofNhnOBsshLgvLG0W_ThlY-L5i6FAsTDp9fWKs_hS4VMn1cb8iexi3Oljcy7255J-wWRSaAMcm4KzZNc3kS_HR7NdfGlu9zmjE22rnmlZS60OEvjhqU-SKJBsalHAiFbAWTemHuk5jlB7P2sFiM4JAxIuznq23s0WrNM0oQTRi6xb0bMglGuBmyvPkoh1jMAGklHStprNoxwY_S2aKiUA",
		Server:    "tcp://192.168.14.120:8055", // 127.0.0.1:1883
		Identifer: "connect",
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

	select {}

	m.UnSubDeviceControl()

	time.Sleep(3 * time.Second)
}

// SendMessageToSDK 用于 deviceControl 的测试
func SendMessageToSDK() {
	client := &http.Client{}

	params := `
		{
			"params":{
				"reply":"this is control test"
			},
			"thing_id":"iott-yVAwx9rb8j"
		}
	`

	requst, err := http.NewRequest("POST", "http://iot-api.qingcloud.com:8889/api/v1/devices/iotd-59f685ce-70f9-4485-9985-271d5dfb9475/call/connect", strings.NewReader(params))
	if err != nil {
		fmt.Println("NewRequest err:", err.Error())
		return
	}
	requst.Header.Set("Authorization", " QCUUCDRNSJECJRCMOMPPHP:signaturea")

	resp, err := client.Do(requst)
	if err != nil {
		fmt.Println("Do err:", err.Error())
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	fmt.Println(string(body))
}
