package main

import (
	"context"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"

	"git.internal.yunify.com/iot-sdk/device-sdk-go/internal/constant"
	"git.internal.yunify.com/iot-sdk/device-sdk-go/internal/define"

	"git.internal.yunify.com/iot-sdk/device-sdk-go/internal/register"
	"git.internal.yunify.com/iot-sdk/device-sdk-go/mqtt"
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

	connect            bool // 上线
	pubProperty        bool // 上报属性
	pubPropertyByMqtts bool // 上报属性
	pubEvent           bool // 上报事件
	serviceContol      bool // 服务调用
	all                bool // 上线、上报属性、上报事件、服务调用

	reg           bool // 动态注册
	regAndConnect bool // 动态注册并上线设备
)

func init() {

	// 通过命令行参数运行不同功能
	flag.StringVar(&configPath, "conf", "./config.yml", "")
	flag.BoolVar(&connect, "c", false, "")
	flag.BoolVar(&pubProperty, "p", false, "")
	flag.BoolVar(&pubPropertyByMqtts, "ps", false, "")
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

	if pubPropertyByMqtts {
		PubPropertyFuncByMQTTS()
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
		DynamicRegistryAndConnect()
	}

	select {}
}

// ConnectFunc 提供设备上线功能
func ConnectFunc() {
	options := &mqtt.Options{
		Token:           conf.Device.Token,
		AutoReconnect:   conf.Device.AutoReconnect,
		LostConnectChan: make(chan bool),
		Server:          conf.Mqttbroker.AddressMqtt,
	}

	m, err := mqtt.InitWithToken(options)
	if err != nil {
		panic(err)
	}

	// 连接
	err = m.Connect()
	if err != nil {
		panic(err)
	}

	// 掉线后的处理动作
	go func(o *mqtt.Options) {
		for {
			select {
			case <-options.LostConnectChan:
				// 如果不重连，则退出程序
				if !o.AutoReconnect {
					fmt.Println("not reconnect to ehub/ihub, procedure will quit!")
					os.Exit(0)
					return
				}
				// 重连，则提示目前暂时掉线
				fmt.Println("lost connect to ehub/ihub, will auto reconnect!")
			}
		}
	}(options)
}

// PubPropertyFunc 在 0 ～ 100 范围内上报温度属性值
func PubPropertyFunc() {
	options := &mqtt.Options{
		Token:           conf.Device.Token,
		AutoReconnect:   conf.Device.AutoReconnect,
		LostConnectChan: make(chan bool),
		Server:          conf.Mqttbroker.AddressMqtt,
		PropertyType:    constant.PROPERTY_TYPE_BASE,
	}

	m, err := mqtt.InitWithToken(options)
	if err != nil {
		panic(err)
	}

	// 连接
	err = m.Connect()
	if err != nil {
		panic(err)
	}

	// 失去连接后的处理动作
	go func(o *mqtt.Options) {
		for {
			select {
			case <-options.LostConnectChan:
				// 如果不重连，则退出程序
				if !o.AutoReconnect {
					fmt.Println("not reconnect to ehub/ihub, procedure will quit!")
					os.Exit(0)
					return
				}
				// 重连，则提示目前暂时掉线
				fmt.Println("lost connect to ehub/ihub, will auto reconnect!")
			}
		}
	}(options)

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
		fmt.Println("DeviceTemprature:", DeviceTemprature)
		time.Sleep(2 * time.Second)
		DeviceTemprature++
		if DeviceTemprature < 0 || DeviceTemprature > 100 {
			DeviceTemprature = float64(rand.Int63n(int64(HIGH) - int64(LOW)))
		}
		data["temp"] = DeviceTemprature
	}
}

// PubPropertyFuncByMQTTS 通过 mqtts 在 0 ～ 100 范围内上报温度属性值
func PubPropertyFuncByMQTTS() {
	options := &mqtt.Options{
		Token:           conf.Device.Token,
		AutoReconnect:   conf.Device.AutoReconnect,
		LostConnectChan: make(chan bool),
		Server:          conf.Mqttbroker.AddressMqtts,
		PropertyType:    constant.PROPERTY_TYPE_BASE,
		// 如果提供证书路径，将会使用 mqtts 进行通信
		CertFilePath: "cert/iot.qingcloud.com.cer",
	}

	m, err := mqtt.InitWithToken(options)
	if err != nil {
		panic(err)
	}

	// 连接
	err = m.Connect()
	if err != nil {
		panic(err)
	}

	// 失去连接后的处理动作
	go func(o *mqtt.Options) {
		for {
			select {
			case <-options.LostConnectChan:
				// 如果不重连，则退出程序
				if !o.AutoReconnect {
					fmt.Println("not reconnect to ehub/ihub, procedure will quit!")
					os.Exit(0)
					return
				}
				// 重连，则提示目前暂时掉线
				fmt.Println("lost connect to ehub/ihub, will auto reconnect!")
			}
		}
	}(options)

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
		fmt.Println("DeviceTemprature:", DeviceTemprature)
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
	options := &mqtt.Options{
		Token:           conf.Device.Token,
		AutoReconnect:   conf.Device.AutoReconnect,
		LostConnectChan: make(chan bool),
		Server:          conf.Mqttbroker.AddressMqtt,
		PropertyType:    constant.PROPERTY_TYPE_BASE,
	}

	m, err := mqtt.InitWithToken(options)
	if err != nil {
		panic(err)
	}

	// 连接
	err = m.Connect()
	if err != nil {
		panic(err)
	}

	// 失去连接后的处理动作
	go func(o *mqtt.Options) {
		for {
			select {
			case <-options.LostConnectChan:
				// 如果不重连，则退出程序
				if !o.AutoReconnect {
					fmt.Println("not reconnect to ehub/ihub, procedure will quit!")
					os.Exit(0)
					return
				}
				// 重连，则提示目前暂时掉线
				fmt.Println("lost connect to ehub/ihub, will auto reconnect!")
			}
		}
	}(options)

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
			ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
			reply, err := m.PubEvent(ctx, eventData, eventIdentifier)
			if err != nil {
				panic(err)
			}
			fmt.Printf("PubEvent reply:%+v\n", reply)
		}
		time.Sleep(2 * time.Second)
	}
}

// ServiceDeviceControlFunc 服务调用
func ServiceDeviceControlFunc() {

	serviceIdentifer := "setTemperature" // 服务调用的 服务 identifer
	inputIdentifier := "temperature"     // 执行服务调用改变的参数值

	options := &mqtt.Options{
		Token:           conf.Device.Token,
		AutoReconnect:   conf.Device.AutoReconnect,
		LostConnectChan: make(chan bool),
		Server:          conf.Mqttbroker.AddressMqtt,
		PropertyType:    constant.PROPERTY_TYPE_BASE,

		DeviceHandlers: []mqtt.DeviceControlHandler{
			mqtt.DeviceControlHandler{
				ServiceIdentifer: serviceIdentifer,
				InputIdentifier:  inputIdentifier,
				ServiceHandler:   DeviceControlCallback,
			},
		},
	}

	m, err := mqtt.InitWithToken(options)
	if err != nil {
		panic(err)
	}

	// 连接
	err = m.Connect()
	if err != nil {
		panic(err)
	}

	// 失去连接后的处理动作
	go func(o *mqtt.Options) {
		for {
			select {
			case <-options.LostConnectChan:
				// 如果不重连，则退出程序
				if !o.AutoReconnect {
					fmt.Println("not reconnect to ehub/ihub, procedure will quit!")
					os.Exit(0)
					return
				}
				// 重连，则提示目前暂时掉线
				fmt.Println("lost connect to ehub/ihub, will auto reconnect!")
			}
		}
	}(options)

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
			fmt.Println("DeviceTemprature:", DeviceTemprature)

			DeviceTemprature++
			if DeviceTemprature < 0 || DeviceTemprature > 100 {
				DeviceTemprature = float64(rand.Int63n(int64(HIGH) - int64(LOW)))
			}
			data["temp"] = DeviceTemprature
			time.Sleep(2 * time.Second)
		}
	}()

	// 服务调用
	m.SubDeviceControl(serviceIdentifer)

}

// PropertyAndEventAndServiceFunc 提供全功能 demo
func PropertyAndEventAndServiceFunc() {

	eventIdentifier := "temperatureEvent" // 上报事件的 事件 identifer
	serviceIdentifer := "setTemperature"  // 服务调用的 服务 identifer
	inputIdentifier := "temperature"      // 执行服务调用改变的参数值

	options := &mqtt.Options{
		Token:           conf.Device.Token,
		AutoReconnect:   conf.Device.AutoReconnect,
		LostConnectChan: make(chan bool),
		Server:          conf.Mqttbroker.AddressMqtt,
		PropertyType:    constant.PROPERTY_TYPE_BASE,

		DeviceHandlers: []mqtt.DeviceControlHandler{
			mqtt.DeviceControlHandler{
				ServiceIdentifer: serviceIdentifer,
				InputIdentifier:  inputIdentifier,
				ServiceHandler:   DeviceControlCallback,
			},
		},
	}

	m, err := mqtt.InitWithToken(options)
	if err != nil {
		panic(err)
	}

	// 连接
	err = m.Connect()
	if err != nil {
		panic(err)
	}

	// 失去连接后的处理动作
	go func(o *mqtt.Options) {
		for {
			select {
			case <-options.LostConnectChan:
				// 如果不重连，则退出程序
				if !o.AutoReconnect {
					fmt.Println("not reconnect to ehub/ihub, procedure will quit!")
					os.Exit(0)
					return
				}
				// 重连，则提示目前暂时掉线
				fmt.Println("lost connect to ehub/ihub, will auto reconnect!")
			}
		}
	}(options)

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
			fmt.Println("DeviceTemprature:", DeviceTemprature)

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

	// 服务调用
	m.SubDeviceControl(serviceIdentifer)
}

// DeviceControlCallback 服务调用的回调函数
func DeviceControlCallback(inputIdentifier string, msg *define.Message) error {
	for k, v := range msg.Params {
		if k == inputIdentifier {

			// 将设备温度调节为服务下发的温度值
			// float64 为 input 对应的类型
			DeviceTemprature = v.(float64)
		}
	}
	return nil
}

// -------------------------

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

	options := &mqtt.Options{
		MiddleCredential:       conf.Registry.MiddleCredential,
		DynamocRegisterAddress: conf.Registry.ServiceAddress,
		AutoReconnect:          conf.Device.AutoReconnect,
		LostConnectChan:        make(chan bool),
		Server:                 conf.Mqttbroker.AddressMqtt,
		PropertyType:           constant.PROPERTY_TYPE_BASE,
	}
	m, err := mqtt.InitWithMiddleCredential(options)
	if err != nil {
		panic(err)
	}

	// 连接
	err = m.Connect()
	if err != nil {
		panic(err)
	}

	// 掉线后的处理动作
	go func(o *mqtt.Options) {
		for {
			select {
			case <-options.LostConnectChan:
				// 如果不重连，则退出程序
				if !o.AutoReconnect {
					fmt.Println("not reconnect to ehub/ihub, procedure will quit!")
					os.Exit(0)
					return
				}
				// 重连，则提示目前暂时掉线
				fmt.Println("lost connect to ehub/ihub, will auto reconnect!")
			}
		}
	}(options)
}
