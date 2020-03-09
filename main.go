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
	{
		go deviceControl()
		time.Sleep(time.Second * 5)
		SendMessageToSDK()
	}
}

// 设备连接
func deviceConnect() {
	options := &index.Options{
		Token:  "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY3IiOiIxIiwiYXVkIjoiaWFtIiwiYXpwIjoiaWFtIiwiY3VpZCI6ImlhbXItajI3ZXAzZmciLCJlaXNrIjoiVXd1YXktY0s2X2xiTUdwcXJmaTNoQlk3anZoTlA4N0NCeHRjN1BLbzYwdz0iLCJleHAiOjE2MDk0ODQxOTQsImlhdCI6MTU3Nzk0ODE5NCwiaXNzIjoic3RzIiwianRpIjoiVWpJNFdQQW9wNWdQNldPdHJIUUlBVSIsIm5iZiI6MCwib3JnaSI6ImlvdGQtNGQ1NTUyZTAtYWUyNy00OTc1LTllMmEtYjk2NTRhZjI1NjM2Iiwib3d1ciI6InVzci1rZUF5dG16MSIsInByZWYiOiJxcm46cWluZ2Nsb3VkOmlhbToiLCJydHlwIjoicm9sZSIsInN1YiI6InN0cyIsInRoaWQiOiJpb3R0LXlWQXd4OXJiOGoiLCJ0eXAiOiJJRCJ9.NDId6MS_Fi-9mCuUaBeS4sufhoWPCihz5TSgyscD1LMdvSs6KKXaND2fmDhlJcFi3-nbTZS32LR_fx8cYS8_8pHNF2pdyfXStYsm1sbBg6G7mfCXmXLywVfzUUxSgJbXJ7Px1oIIPjcuPCmlEK4BtDyK5a5Ncxw9NO0aZxKviNqPKMOqQAPP8_2Ev6MGQ4SwsLuZP3dE75bTp02XID1xCGY_0ABIPhHQrypqs2T-_h1DE-5MZegSL5sUjjgha4AVH_2xzPcgLKO709e77tWhu5BpJXUmUfTlZwUp3PoDG4eNYC3gqVEgAkZtUxjoCvGXypqV7lV8YudYmrN7BBuXmw",
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
		Token:        "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY3IiOiIxIiwiYXVkIjoiaWFtIiwiYXpwIjoiaWFtIiwiY3VpZCI6ImlhbXItajI3ZXAzZmciLCJlaXNrIjoiVXd1YXktY0s2X2xiTUdwcXJmaTNoQlk3anZoTlA4N0NCeHRjN1BLbzYwdz0iLCJleHAiOjE2MDk0ODQxOTQsImlhdCI6MTU3Nzk0ODE5NCwiaXNzIjoic3RzIiwianRpIjoiVWpJNFdQQW9wNWdQNldPdHJIUUlBVSIsIm5iZiI6MCwib3JnaSI6ImlvdGQtNGQ1NTUyZTAtYWUyNy00OTc1LTllMmEtYjk2NTRhZjI1NjM2Iiwib3d1ciI6InVzci1rZUF5dG16MSIsInByZWYiOiJxcm46cWluZ2Nsb3VkOmlhbToiLCJydHlwIjoicm9sZSIsInN1YiI6InN0cyIsInRoaWQiOiJpb3R0LXlWQXd4OXJiOGoiLCJ0eXAiOiJJRCJ9.NDId6MS_Fi-9mCuUaBeS4sufhoWPCihz5TSgyscD1LMdvSs6KKXaND2fmDhlJcFi3-nbTZS32LR_fx8cYS8_8pHNF2pdyfXStYsm1sbBg6G7mfCXmXLywVfzUUxSgJbXJ7Px1oIIPjcuPCmlEK4BtDyK5a5Ncxw9NO0aZxKviNqPKMOqQAPP8_2Ev6MGQ4SwsLuZP3dE75bTp02XID1xCGY_0ABIPhHQrypqs2T-_h1DE-5MZegSL5sUjjgha4AVH_2xzPcgLKO709e77tWhu5BpJXUmUfTlZwUp3PoDG4eNYC3gqVEgAkZtUxjoCvGXypqV7lV8YudYmrN7BBuXmw",
		Server:       "tcp://192.168.14.120:8055", // 127.0.0.1:1883 192.168.14.120:1883
		PropertyType: index.PROPERTY_TYPE_BASE,
		MessageID:    "message-device.1",
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

	data := index.PropertyKV{
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
		Token:           "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY3IiOiIxIiwiYXVkIjoiaWFtIiwiYXpwIjoiaWFtIiwiY3VpZCI6ImlhbXItajI3ZXAzZmciLCJlaXNrIjoiVXd1YXktY0s2X2xiTUdwcXJmaTNoQlk3anZoTlA4N0NCeHRjN1BLbzYwdz0iLCJleHAiOjE2MDk0ODQxOTQsImlhdCI6MTU3Nzk0ODE5NCwiaXNzIjoic3RzIiwianRpIjoiVWpJNFdQQW9wNWdQNldPdHJIUUlBVSIsIm5iZiI6MCwib3JnaSI6ImlvdGQtNGQ1NTUyZTAtYWUyNy00OTc1LTllMmEtYjk2NTRhZjI1NjM2Iiwib3d1ciI6InVzci1rZUF5dG16MSIsInByZWYiOiJxcm46cWluZ2Nsb3VkOmlhbToiLCJydHlwIjoicm9sZSIsInN1YiI6InN0cyIsInRoaWQiOiJpb3R0LXlWQXd4OXJiOGoiLCJ0eXAiOiJJRCJ9.NDId6MS_Fi-9mCuUaBeS4sufhoWPCihz5TSgyscD1LMdvSs6KKXaND2fmDhlJcFi3-nbTZS32LR_fx8cYS8_8pHNF2pdyfXStYsm1sbBg6G7mfCXmXLywVfzUUxSgJbXJ7Px1oIIPjcuPCmlEK4BtDyK5a5Ncxw9NO0aZxKviNqPKMOqQAPP8_2Ev6MGQ4SwsLuZP3dE75bTp02XID1xCGY_0ABIPhHQrypqs2T-_h1DE-5MZegSL5sUjjgha4AVH_2xzPcgLKO709e77tWhu5BpJXUmUfTlZwUp3PoDG4eNYC3gqVEgAkZtUxjoCvGXypqV7lV8YudYmrN7BBuXmw",
		Server:          "tcp://192.168.14.120:8055", // 127.0.0.1:1883  192.168.14.120:8055
		EventIdentifier: "serviceStatus",
		MessageID:       "message-device.1",
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

	// output
	data := index.PropertyKV{
		"ServiceName":   "qqq",
		"ServiceStatus": "aaa",
	}
	reply, err := m.PubEvent(context.Background(), data)
	if err != nil {
		panic(err)
	}
	fmt.Printf("PubEvent reply:%+v\n", reply)
	select {}
}

func deviceControl() {
	options := &index.Options{
		Token:     "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY3IiOiIxIiwiYXVkIjoiaWFtIiwiYXpwIjoiaWFtIiwiY3VpZCI6ImlhbXItajI3ZXAzZmciLCJlaXNrIjoiVXd1YXktY0s2X2xiTUdwcXJmaTNoQlk3anZoTlA4N0NCeHRjN1BLbzYwdz0iLCJleHAiOjE2MDk0ODQxOTQsImlhdCI6MTU3Nzk0ODE5NCwiaXNzIjoic3RzIiwianRpIjoiVWpJNFdQQW9wNWdQNldPdHJIUUlBVSIsIm5iZiI6MCwib3JnaSI6ImlvdGQtNGQ1NTUyZTAtYWUyNy00OTc1LTllMmEtYjk2NTRhZjI1NjM2Iiwib3d1ciI6InVzci1rZUF5dG16MSIsInByZWYiOiJxcm46cWluZ2Nsb3VkOmlhbToiLCJydHlwIjoicm9sZSIsInN1YiI6InN0cyIsInRoaWQiOiJpb3R0LXlWQXd4OXJiOGoiLCJ0eXAiOiJJRCJ9.NDId6MS_Fi-9mCuUaBeS4sufhoWPCihz5TSgyscD1LMdvSs6KKXaND2fmDhlJcFi3-nbTZS32LR_fx8cYS8_8pHNF2pdyfXStYsm1sbBg6G7mfCXmXLywVfzUUxSgJbXJ7Px1oIIPjcuPCmlEK4BtDyK5a5Ncxw9NO0aZxKviNqPKMOqQAPP8_2Ev6MGQ4SwsLuZP3dE75bTp02XID1xCGY_0ABIPhHQrypqs2T-_h1DE-5MZegSL5sUjjgha4AVH_2xzPcgLKO709e77tWhu5BpJXUmUfTlZwUp3PoDG4eNYC3gqVEgAkZtUxjoCvGXypqV7lV8YudYmrN7BBuXmw",
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

	requst, err := http.NewRequest("POST", "http://iot-api.qingcloud.com:8889/api/v1/devices/iotd-4d5552e0-ae27-4975-9e2a-b9654af25636/call/connect", strings.NewReader(params))
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
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
}
