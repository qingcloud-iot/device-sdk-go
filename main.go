package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"time"

	"git.internal.yunify.com/iot-sdk/device-sdk-go/index"
	"git.internal.yunify.com/iot-sdk/device-sdk-go/mqtt"
	mqttp "github.com/eclipse/paho.mqtt.golang"
)

const (
	LOW  float64 = 30 // 设备温度下限
	HIGH float64 = 50 // 设备温度上限
)

var VALUE float64 = 30 // 模拟设备值

var (
	configPath string

	connect       bool // 上线
	pubProperty   bool // 上报属性
	pubEvent      bool // 上报事件
	serviceContol bool // 设备控制
	all           bool // 上线、上报属性、上报事件、设备控制
)

func init() {
	// 加载配置文件

	// 通过命令行参数运行不同功能
	flag.StringVar(&configPath, "conf", "./config.yml", "")
	flag.BoolVar(&connect, "c", false, "")
	flag.BoolVar(&pubProperty, "p", false, "")
	flag.BoolVar(&pubEvent, "e", false, "")
	flag.BoolVar(&serviceContol, "s", false, "")
	flag.BoolVar(&all, "all", false, "")
	flag.Parse()

	InitConfig()
}

func main() {

	if connect {
		ConnectFunc()
	}
	if pubProperty {
		PubPropertyFunc()
	}

	if pubEvent {
		PubEventFunc()
	}

	if serviceContol {
		ServiceDeviceControlFunc()
	}

	if all {
		PropertyAndEventAndServiceFunc()
	}
}

// ConnectFunc 提供设备上线功能
func ConnectFunc() {
	token := conf.Device.Token
	entityID, modelID, err := mqtt.ParseToken(token)
	if err != nil {
		panic("Parse token error: " + err.Error())
	}

	options := &mqtt.Options{
		Token:     token,
		Server:    conf.Mqttbroker.Addr, // 127.0.0.1:1883 192.168.14.120:1883
		MessageID: "message-device.1",
		EntityId:  entityID,
		ModelId:   modelID,
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

// PubPropertyFunc 在 0 ～ 100 范围内上报温度属性值
func PubPropertyFunc() {
	token := conf.Device.Token
	entityID, modelID, err := mqtt.ParseToken(token)
	if err != nil {
		panic("Parse token error: " + err.Error())
	}

	options := &mqtt.Options{
		Token:        token,
		Server:       conf.Mqttbroker.Addr, // 127.0.0.1:1883 192.168.14.120:1883
		MessageID:    "message-device.1",
		PropertyType: index.PROPERTY_TYPE_BASE,
		EntityId:     entityID,
		ModelId:      modelID,
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
		"temp": VALUE,
	}

	// 上报属性
	for {
		ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
		_, err := m.PubProperty(ctx, data)
		if err != nil {
			panic(err)
		}
		// fmt.Println("PubPropertySync reply", reply)
		time.Sleep(2 * time.Second)
		VALUE++
		if VALUE < 0 || VALUE > 100 {
			VALUE = float64(rand.Int63n(int64(HIGH) - int64(LOW)))
		}
		data["temp"] = VALUE
	}
}

// PubEventFunc 上报事件
func PubEventFunc() {
	token := conf.Device.Token
	entityID, modelID, err := mqtt.ParseToken(token)
	if err != nil {
		panic("Parse token error: " + err.Error())
	}

	options := &mqtt.Options{
		Token:     token,
		Server:    conf.Mqttbroker.Addr, // 127.0.0.1:1883 192.168.14.120:1883
		MessageID: "message-device.1",
		EntityId:  entityID,
		ModelId:   modelID,
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

	eventIdentifier := "temperatureEvent" // 上报事件的 事件 identifer

	// 上报事件
	for {
		VALUE = float64(rand.Intn(50)) + LOW - 15

		// 当温度低于 30 超过 50 时，上报事件
		if VALUE < 30 || VALUE > 50 {
			eventData := index.PropertyKV{
				"temperature": float64(VALUE),
				"region":      "北京",
				"name":        "锅炉",
			}
			reply, err := m.PubEvent(context.Background(), eventData, eventIdentifier)
			if err != nil {
				panic(err)
			}
			fmt.Printf("PubEvent reply:%+v\n", reply)
		}
		time.Sleep(2 * time.Second)
	}
}

// ServiceDeviceControlFunc 设备控制
func ServiceDeviceControlFunc() {
	token := conf.Device.Token

	entityID, modelID, err := mqtt.ParseToken(token)
	if err != nil {
		panic("Parse token error: " + err.Error())
	}
	serviceIdentifer := "setTemperature" // 服务调用的 服务 identifer

	options := &mqtt.Options{
		Token:        token,
		Server:       conf.Mqttbroker.Addr, // 127.0.0.1:1883 192.168.14.120:1883
		PropertyType: index.PROPERTY_TYPE_BASE,
		MessageID:    "message-device.1",
		EntityId:     entityID,
		ModelId:      modelID,
	}

	// 供设备控制使用
	options.SetDeviceControlHandler(func(client mqttp.Client, msg mqttp.Message) {
		fmt.Printf("[sdk-go sub] topic:%s, message:%s\n", msg.Topic(), string(msg.Payload()))
		switch {
		case msg.Topic() == mqtt.BuildServiceControlReply(modelID, entityID, serviceIdentifer):
			RecvDeviceControlReply(client, msg)
		default:
		}
	})

	m, err := mqtt.NewMqtt(options)
	if err != nil {
		panic(err)
	}

	// 连接
	err = m.Connect()
	if err != nil {
		panic(err)
	}

	// 上报属性
	go func() {
		for {
			data := index.PropertyKV{
				"temp": VALUE,
			}
			ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
			_, err := m.PubProperty(ctx, data)
			if err != nil {
				panic(err)
			}
			// fmt.Println("PubPropertySync reply", reply)
			VALUE++
			if VALUE < 0 || VALUE > 100 {
				VALUE = float64(rand.Int63n(int64(HIGH) - int64(LOW)))
			}
			data["temp"] = VALUE
			time.Sleep(2 * time.Second)
		}
	}()

	// 设备控制
	m.SubDeviceControl(serviceIdentifer)

}

// PropertyAndEventAndServiceFunc 提供全功能 demo
func PropertyAndEventAndServiceFunc() {
	token := conf.Device.Token

	entityID, modelID, err := mqtt.ParseToken(token)
	if err != nil {
		panic("Parse token error: " + err.Error())
	}

	eventIdentifier := "temperatureEvent" // 上报事件的 事件 identifer
	serviceIdentifer := "setTemperature"  // 服务调用的 服务 identifer

	options := &mqtt.Options{
		Token:        token,
		Server:       conf.Mqttbroker.Addr, // 127.0.0.1:1883 192.168.14.120:1883
		PropertyType: index.PROPERTY_TYPE_BASE,
		MessageID:    "message-device.1",
		EntityId:     entityID,
		ModelId:      modelID,
	}

	// 供设备控制使用的回调函数
	options.SetDeviceControlHandler(func(client mqttp.Client, msg mqttp.Message) {
		fmt.Printf("[sdk-go sub] topic:%s, message:%s\n", msg.Topic(), string(msg.Payload()))
		switch {
		case msg.Topic() == mqtt.BuildServiceControlReply(modelID, entityID, serviceIdentifer):
			RecvDeviceControlReply(client, msg)
		default:
		}
	})

	m, err := mqtt.NewMqtt(options)
	if err != nil {
		panic(err)
	}

	// 连接
	err = m.Connect()
	if err != nil {
		panic(err)
	}

	// 上报属性
	go func() {
		for {
			data := index.PropertyKV{
				"temp": VALUE,
			}
			ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
			_, err := m.PubProperty(ctx, data)
			if err != nil {
				panic(err)
			}
			// fmt.Println("PubPropertySync reply", reply)
			VALUE++
			if VALUE < 0 || VALUE > 100 {
				VALUE = float64(rand.Int63n(int64(HIGH) - int64(LOW)))
			}
			data["temp"] = VALUE
			time.Sleep(2 * time.Second)
		}
	}()

	// 上报事件
	go func() {
		for {
			// 当温度低于 30 超过 50 时，上报事件
			if VALUE < 30 || VALUE > 50 {
				eventData := index.PropertyKV{
					"temperature": float64(VALUE),
					"region":      "北京",
					"name":        "锅炉",
				}
				reply, err := m.PubEvent(context.Background(), eventData, eventIdentifier)
				if err != nil {
					panic(err)
				}
				fmt.Printf("PubEvent reply:%+v\n", reply)
			}
			time.Sleep(2 * time.Second)
		}
	}()

	// 设备控制
	m.SubDeviceControl(serviceIdentifer)
}

// recvDeviceControlReply 订阅消息后的回调函数，实现具体业务逻辑
func RecvDeviceControlReply(client mqtt.Client, msg mqtt.Message) {

	topic := msg.Topic()
	payload := msg.Payload()

	fmt.Printf("[recvDeviceControlReply] topic:%s payload:%s\n", topic, string(payload))
	message, err := mqtt.ParseMessage(payload)
	if err != nil {
		fmt.Printf("recvDeviceControlReply err:%s", err.Error())
		return
	}
	VALUE = message.Params["temperature"].(float64)

	reply := &index.Reply{
		Id:   message.Id,
		Code: index.RPC_SUCCESS,
		Data: make(index.PropertyKV),
	}

	reply.Data = message.Params

	data, err := json.Marshal(reply)
	if err != nil {
		fmt.Printf("[recvDeviceControlReply] err:%s\n", err.Error())
		return
	}
	/*
		{
		    "code":200,
		    "id":"49be65b8-2746-41e8-b314-afd724f2213e",
		    "data":{
		        "temperature":30
		    }
		}
	*/
	token := client.Publish(topic+"_reply", byte(0), false, data)
	if token.Error() != nil {
		fmt.Printf("[recvDeviceControlReply] err:%s\n", err.Error())
	} else {
		fmt.Printf("[recvDeviceControlReply] success\n")
	}
}
