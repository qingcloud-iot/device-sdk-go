package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"git.internal.yunify.com/iot-sdk/device-sdk-go/internal/constant"
	"git.internal.yunify.com/iot-sdk/device-sdk-go/internal/define"
	"math/rand"
	"time"

	"git.internal.yunify.com/iot-sdk/device-sdk-go/mqtt"
	"git.internal.yunify.com/iot-sdk/device-sdk-go/register"
	mqttp "github.com/eclipse/paho.mqtt.golang"
	uuid "github.com/satori/go.uuid"
)

const (
	// LOW 设备温度下限
	LOW float64 = 30
	// HIGH 设备温度上限
	HIGH float64 = 50
)

// DeviceTemprature 模拟设备温度值
var DeviceTemprature float64 = 30

var (
	configPath string

	connect       bool // 上线
	pubProperty   bool // 上报属性
	pubEvent      bool // 上报事件
	serviceContol bool // 设备控制
	all           bool // 上线、上报属性、上报事件、设备控制

	reg           bool // 动态注册
	regAndConnect bool // 动态注册并上线设备
)

func init() {

	// 通过命令行参数运行不同功能
	flag.StringVar(&configPath, "conf", "./config.yml", "")
	flag.BoolVar(&connect, "c", false, "")
	flag.BoolVar(&pubProperty, "p", false, "")
	flag.BoolVar(&pubEvent, "e", false, "")
	flag.BoolVar(&serviceContol, "s", false, "")
	flag.BoolVar(&all, "all", false, "")
	flag.BoolVar(&reg, "r", false, "")
	flag.BoolVar(&regAndConnect, "rc", false, "")

	flag.Parse()

	// 加载配置文件
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

	if reg {
		DynamicRegistry()
	}

	if regAndConnect {
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
		Server:    conf.Mqttbroker.Address,
		MessageID: uuid.NewV4().String(),
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
		Server:       conf.Mqttbroker.Address,
		MessageID:    uuid.NewV4().String(),
		PropertyType: constant.PROPERTY_TYPE_BASE,
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

	data := define.PropertyKV{
		"temp": DeviceTemprature,
	}

	// 上报属性
	for {
		ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
		_, err := m.PubProperty(ctx, data)
		if err != nil {
			panic(err)
		}
		time.Sleep(2 * time.Second)
		DeviceTemprature++
		if DeviceTemprature < 0 || DeviceTemprature > 100 {
			DeviceTemprature = float64(rand.Int63n(int64(HIGH) - int64(LOW)))
		}
		data["temp"] = DeviceTemprature
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
		Server:    conf.Mqttbroker.Address,
		MessageID: uuid.NewV4().String(),
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
		DeviceTemprature = float64(rand.Intn(50)) + LOW - 15

		// 当温度低于 30 超过 50 时，上报事件
		if DeviceTemprature < 30 || DeviceTemprature > 50 {
			var reason int
			if DeviceTemprature < 30 {
				reason = 0
			}
			if DeviceTemprature > 50 {
				reason = 1
			}

			eventData := define.PropertyKV{
				"temperature": float64(DeviceTemprature),
				"reason":      reason,
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
		Server:       conf.Mqttbroker.Address,
		PropertyType: constant.PROPERTY_TYPE_BASE,
		MessageID:    uuid.NewV4().String(),
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
			data := define.PropertyKV{
				"temp": DeviceTemprature,
			}
			ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
			_, err := m.PubProperty(ctx, data)
			if err != nil {
				panic(err)
			}

			DeviceTemprature++
			if DeviceTemprature < 0 || DeviceTemprature > 100 {
				DeviceTemprature = float64(rand.Int63n(int64(HIGH) - int64(LOW)))
			}
			data["temp"] = DeviceTemprature
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
		Server:       conf.Mqttbroker.Address,
		PropertyType: constant.PROPERTY_TYPE_BASE,
		MessageID:    uuid.NewV4().String(),
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
			data := define.PropertyKV{
				"temp": DeviceTemprature,
			}
			ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
			_, err := m.PubProperty(ctx, data)
			if err != nil {
				panic(err)
			}

			DeviceTemprature++
			if DeviceTemprature < 0 || DeviceTemprature > 100 {
				DeviceTemprature = float64(rand.Int63n(int64(HIGH) - int64(LOW)))
			}
			data["temp"] = DeviceTemprature
			time.Sleep(2 * time.Second)
		}
	}()

	// 上报事件
	go func() {
		for {
			// 当温度低于 30 超过 50 时，上报事件
			if DeviceTemprature < 30 || DeviceTemprature > 50 {
				var reason int
				if DeviceTemprature < 30 {
					reason = 0 // 温度过低
				}
				if DeviceTemprature > 50 {
					reason = 1 // 温度过高
				}

				eventData := define.PropertyKV{
					"temperature": float64(DeviceTemprature),
					"reason":      reason,
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

// RecvDeviceControlReply 订阅消息后的回调函数，实现具体业务逻辑
func RecvDeviceControlReply(client mqtt.Client, msg mqtt.Message) {

	topic := msg.Topic()
	payload := msg.Payload()

	fmt.Printf("[recvDeviceControlReply] topic:%s payload:%s\n", topic, string(payload))
	message, err := mqtt.ParseMessage(payload)
	if err != nil {
		fmt.Printf("recvDeviceControlReply err:%s", err.Error())
		return
	}

	// 将设备温度调节为服务下发的温度值
	DeviceTemprature = message.Params["temperature"].(float64)

	reply := &define.Reply{
		ID:   message.ID,
		Code: constant.RPC_SUCCESS,
		Data: make(define.PropertyKV),
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

// DynamicRegistry 设备的动态注册
func DynamicRegistry() {
	midCredential := conf.Registry.MiddleCredential

	r := register.NewRegister(conf.Registry.ServiceAddress)
	resp, err := r.DynamicRegistry(midCredential)
	if err != nil {
		fmt.Printf("%s dynamic registry failed, error: %s\n", midCredential, err.Error())
		return
	}
	fmt.Printf("%s dynamic registry success, ID:%s, device_name:%s, token:%s\n", midCredential, resp.ID, resp.DeviceName, resp.Token)
}

// DynamicRegistryAndConnect 设备的动态注册并上线
func DynamicRegistryAndConnect() {

}
