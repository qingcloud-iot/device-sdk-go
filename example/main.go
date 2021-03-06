package main

import (
	"context"
	"flag"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/qingcloud-iot/device-sdk-go/constant"
	"github.com/qingcloud-iot/device-sdk-go/define"

	"github.com/qingcloud-iot/device-sdk-go/mqtt"
	"github.com/qingcloud-iot/device-sdk-go/register"
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
	certPath   string

	connect             bool // 上线
	pubProperty         bool // 上报属性
	pubPropertyWithTime bool // 上报属性及自定义时间
	pubPropertyByMqtts  bool // 上报属性
	pubEvent            bool // 上报事件
	serviceContol       bool // 服务调用
	all                 bool // 上线、上报属性、上报事件、服务调用

	reg           bool // 动态注册
	regAndConnect bool // 动态注册并上线设备
	sub           bool // 订阅消息
)

func init() {

	// 通过命令行参数运行不同功能
	flag.StringVar(&configPath, "conf", "./config.yml", "指定配置文件")
	flag.StringVar(&certPath, "ca", "", "tls证书文件(非官方CA机构颁发的证书，请使用本参数)")
	flag.BoolVar(&connect, "c", false, "设备连接")
	flag.BoolVar(&pubProperty, "p", false, "上报属性数据")
	flag.BoolVar(&pubPropertyWithTime, "pt", false, "上报属性数据和时间")
	flag.BoolVar(&pubPropertyByMqtts, "ps", false, "tls加密上报属性数据和时间")
	flag.BoolVar(&pubEvent, "e", false, "上报事件")
	flag.BoolVar(&serviceContol, "s", false, "服务调用")
	flag.BoolVar(&all, "all", false, "")
	flag.BoolVar(&reg, "r", false, "动态注册")
	flag.BoolVar(&regAndConnect, "rc", false, "动态注册并且连接")
	flag.BoolVar(&sub, "sub", false, "订阅设备下行消息")

	flag.Parse()

	// 加载配置文件
	InitConfig()
}

func main() {

	log.SetFlags(log.Lshortfile | log.LstdFlags)

	if connect {
		ConnectFunc()
	}
	if pubProperty {
		PubPropertyFunc()
	}
	if pubPropertyWithTime {
		PubPropertyWithTimeFunc()
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

	if sub {
		Sub()
	}

	select {}
}

// ConnectFunc 提供设备上线功能
func ConnectFunc() {
	options := &mqtt.Options{
		Token:           conf.Device.Token,
		AutoReconnect:   conf.Device.AutoReconnect,
		LostConnectChan: make(chan bool),
		ReConnectChan:   make(chan bool),
		Server:          conf.Mqttbroker.AddressMqtt,

		KeepAlive:            60,               // 心跳间隔, 默认 30s
		MaxReconnectInterval: 20 * time.Minute, // 最大重连间隔, 默认 10 * time.Minute
	}

	m, err := mqtt.InitWithToken(options)
	if err != nil {
		log.Panic(err)
	}

	// 连接
	err = m.Connect()
	if err != nil {
		log.Panic(err)
	}

	// 掉线后的处理动作
	go func(o *mqtt.Options) {
		for {
			select {
			case <-options.LostConnectChan:
				// 如果不重连，则退出程序
				if !o.AutoReconnect {
					log.Println("not reconnect to ehub/ihub, procedure will quit!")
					os.Exit(0)
					return
				}
				// 重连，则提示目前暂时掉线
				log.Println("lost connect to ehub/ihub, will auto reconnect!")
			case ok := <-options.ReConnectChan:
				if ok {
					log.Println("设备已重新连接")
				}
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
		ReConnectChan:   make(chan bool),
		Server:          conf.Mqttbroker.AddressMqtt,
		PropertyType:    constant.PROPERTY_TYPE_BASE,
	}

	m, err := mqtt.InitWithToken(options)
	if err != nil {
		log.Panic(err)
	}

	// 连接
	err = m.Connect()
	if err != nil {
		log.Panic(err)
	}

	// 失去连接后的处理动作
	go func(o *mqtt.Options) {
		for {
			select {
			case <-options.LostConnectChan:
				// 如果不重连，则退出程序
				if !o.AutoReconnect {
					log.Println("not reconnect to ehub/ihub, procedure will quit!")
					os.Exit(0)
					return
				}
				// 重连，则提示目前暂时掉线
				log.Println("lost connect to ehub/ihub, will auto reconnect!")
			case ok := <-options.ReConnectChan:
				if ok {
					log.Println("设备已连接")
				}
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
			log.Panic(err)
		}
		log.Println("DeviceTemprature:", DeviceTemprature)
		time.Sleep(2 * time.Second)
		DeviceTemprature++
		if DeviceTemprature < 0 || DeviceTemprature > 100 {
			DeviceTemprature = float64(rand.Int63n(int64(HIGH) - int64(LOW)))
		}
		data["temp"] = DeviceTemprature
	}
}

// PubPropertyWithTimeFunc 在 0 ～ 100 范围内上报自定义时间及温度属性值
func PubPropertyWithTimeFunc() {
	options := &mqtt.Options{
		Token:           conf.Device.Token,
		AutoReconnect:   conf.Device.AutoReconnect,
		LostConnectChan: make(chan bool),
		ReConnectChan:   make(chan bool),
		Server:          conf.Mqttbroker.AddressMqtt,
		PropertyType:    constant.PROPERTY_TYPE_BASE,
	}

	m, err := mqtt.InitWithToken(options)
	if err != nil {
		log.Panic(err)
	}

	// 连接
	err = m.Connect()
	if err != nil {
		log.Panic(err)
	}

	// 失去连接后的处理动作
	go func(o *mqtt.Options) {
		for {
			select {
			case <-options.LostConnectChan:
				// 如果不重连，则退出程序
				if !o.AutoReconnect {
					log.Println("not reconnect to ehub/ihub, procedure will quit!")
					os.Exit(0)
					return
				}
				// 重连，则提示目前暂时掉线
				log.Println("lost connect to ehub/ihub, will auto reconnect!")
			case ok := <-options.ReConnectChan:
				if ok {
					log.Println("设备已连接")
				}
			}
		}
	}(options)

	data := define.PropertyKVWithTime{
		"temp": &define.PropertyValueAndTime{
			Value: DeviceTemprature,
			Time:  time.Now().UnixNano()/1e6 - 10,
		},
	}

	// 上报属性
	for {
		ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
		_, err := m.PubPropertyWithTime(ctx, data)
		if err != nil {
			log.Panic(err)
		}
		log.Println("DeviceTemprature:", DeviceTemprature)
		time.Sleep(2 * time.Second)
		DeviceTemprature++
		if DeviceTemprature < 0 || DeviceTemprature > 100 {
			DeviceTemprature = float64(rand.Int63n(int64(HIGH) - int64(LOW)))
		}

		data = define.PropertyKVWithTime{
			"temp": &define.PropertyValueAndTime{
				Value: DeviceTemprature,
				Time:  time.Now().UnixNano()/1e6 - 10,
			},
		}
	}
}

// PubPropertyFuncByMQTTS 通过 mqtts 在 0 ～ 100 范围内上报温度属性值
func PubPropertyFuncByMQTTS() {
	options := &mqtt.Options{
		Token:           conf.Device.Token,
		AutoReconnect:   conf.Device.AutoReconnect,
		LostConnectChan: make(chan bool),
		ReConnectChan:   make(chan bool),
		Server:          conf.Mqttbroker.AddressMqtts,
		PropertyType:    constant.PROPERTY_TYPE_BASE,
		// 如果提供证书路径，将会使用 mqtts 进行通信
		CertFilePath: certPath,
		TLS:          true,
	}

	m, err := mqtt.InitWithToken(options)
	if err != nil {
		log.Panic(err)
	}

	// 连接
	err = m.Connect()
	if err != nil {
		log.Panic(err)
	}

	// 失去连接后的处理动作
	go func(o *mqtt.Options) {
		for {
			select {
			case <-options.LostConnectChan:
				// 如果不重连，则退出程序
				if !o.AutoReconnect {
					log.Println("not reconnect to ehub/ihub, procedure will quit!")
					os.Exit(0)
					return
				}
				// 重连，则提示目前暂时掉线
				log.Println("lost connect to ehub/ihub, will auto reconnect!")
			case ok := <-options.ReConnectChan:
				if ok {
					log.Println("设备已连接")
				}
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
			log.Panic(err)
		}
		log.Println("DeviceTemprature:", DeviceTemprature)
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
		ReConnectChan:   make(chan bool),
		Server:          conf.Mqttbroker.AddressMqtt,
		PropertyType:    constant.PROPERTY_TYPE_BASE,
	}

	m, err := mqtt.InitWithToken(options)
	if err != nil {
		log.Panic(err)
	}

	// 连接
	err = m.Connect()
	if err != nil {
		log.Panic(err)
	}

	// 失去连接后的处理动作
	go func(o *mqtt.Options) {
		for {
			select {
			case <-options.LostConnectChan:
				// 如果不重连，则退出程序
				if !o.AutoReconnect {
					log.Println("not reconnect to ehub/ihub, procedure will quit!")
					os.Exit(0)
					return
				}
				// 重连，则提示目前暂时掉线
				log.Println("lost connect to ehub/ihub, will auto reconnect!")
			case ok := <-options.ReConnectChan:
				if ok {
					log.Println("设备已连接")
				}
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
				log.Panic(err)
			}
			log.Printf("PubEvent reply:%+v\n", reply)
		}
		time.Sleep(2 * time.Second)
	}
}

// ServiceDeviceControlFunc 服务调用
func ServiceDeviceControlFunc() {

	serviceIdentifer := "setTemperature" // 服务调用的 服务 identifer

	params := &InAndOutputParameters{
		InputParam1:  "temperature",
		OutputParam1: "result",
		OutputParam2: "temperature",
	}

	options := &mqtt.Options{
		Token:           conf.Device.Token,
		AutoReconnect:   conf.Device.AutoReconnect,
		LostConnectChan: make(chan bool),
		ReConnectChan:   make(chan bool),
		Server:          conf.Mqttbroker.AddressMqtt,
		PropertyType:    constant.PROPERTY_TYPE_BASE,

		DeviceHandlers: []mqtt.DeviceControlHandler{
			mqtt.DeviceControlHandler{
				ServiceIdentifer: serviceIdentifer,
				ServiceHandler:   params,
			},
		},
	}

	m, err := mqtt.InitWithToken(options)
	if err != nil {
		log.Panic(err)
	}

	// 连接
	err = m.Connect()
	if err != nil {
		log.Panic(err)
	}

	// 失去连接后的处理动作
	go func(o *mqtt.Options) {
		for {
			select {
			case <-options.LostConnectChan:
				// 如果不重连，则退出程序
				if !o.AutoReconnect {
					log.Println("not reconnect to ehub/ihub, procedure will quit!")
					os.Exit(0)
					return
				}
				// 重连，则提示目前暂时掉线
				log.Println("lost connect to ehub/ihub, will auto reconnect!")
			case ok := <-options.ReConnectChan:
				if ok {
					log.Println("设备已连接")
				}
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
				log.Panic(err)
			}
			log.Println("DeviceTemprature:", DeviceTemprature)

			DeviceTemprature++
			if DeviceTemprature < 0 || DeviceTemprature > 100 {
				DeviceTemprature = float64(rand.Int63n(int64(HIGH) - int64(LOW)))
			}
			data["temp"] = DeviceTemprature
			time.Sleep(2 * time.Second)
		}
	}()
}

// PropertyAndEventAndServiceFunc 提供全功能 demo
func PropertyAndEventAndServiceFunc() {

	eventIdentifier := "temperatureEvent" // 上报事件的 事件 identifer
	serviceIdentifer := "setTemperature"  // 服务调用的 服务 identifer

	params := &InAndOutputParameters{
		InputParam1:  "temperature",
		OutputParam1: "result",
		OutputParam2: "temperature",
	}

	options := &mqtt.Options{
		Token:           conf.Device.Token,
		AutoReconnect:   conf.Device.AutoReconnect,
		LostConnectChan: make(chan bool),
		ReConnectChan:   make(chan bool),
		Server:          conf.Mqttbroker.AddressMqtt,
		PropertyType:    constant.PROPERTY_TYPE_BASE,

		DeviceHandlers: []mqtt.DeviceControlHandler{
			mqtt.DeviceControlHandler{
				ServiceIdentifer: serviceIdentifer,
				ServiceHandler:   params,
			},
		},
	}

	m, err := mqtt.InitWithToken(options)
	if err != nil {
		log.Panic(err)
	}

	// 连接
	err = m.Connect()
	if err != nil {
		log.Panic(err)
	}

	// 失去连接后的处理动作
	go func(o *mqtt.Options) {
		for {
			select {
			case <-options.LostConnectChan:
				// 如果不重连，则退出程序
				if !o.AutoReconnect {
					log.Println("not reconnect to ehub/ihub, procedure will quit!")
					os.Exit(0)
					return
				}
				// 重连，则提示目前暂时掉线
				log.Println("lost connect to ehub/ihub, will auto reconnect!")
			case ok := <-options.ReConnectChan:
				if ok {
					log.Println("设备已连接")
				}
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
				log.Panic(err)
			}
			log.Println("DeviceTemprature:", DeviceTemprature)

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
					log.Panic(err)
				}
				log.Printf("PubEvent reply:%+v\n", reply)
			}
			time.Sleep(2 * time.Second)
		}
	}()
}

type InAndOutputParameters struct {
	InputParam1  string
	OutputParam1 string
	OutputParam2 string
}

// DeviceControlCallback 服务调用的回调函数
func (p *InAndOutputParameters) Handler(msg *define.Message) define.PropertyKV {

	// 服务调用返回给平台的值 (对应 output 参数)
	callbackResult := make(define.PropertyKV)

	for k, v := range msg.Params {
		// 服务调用调节的值 (对应 input 参数)
		if k == p.InputParam1 {
			// 将设备温度调节为服务下发的温度值, float64 为 input 对应的类型
			// 这里是设置值的相应逻辑，通过设置的成功与否，定义返回值
			assertValue, ok := v.(float64)
			if ok {
				DeviceTemprature = assertValue
				// 如果设置成功
				callbackResult[p.OutputParam1] = 1
				callbackResult[p.OutputParam2] = assertValue
			} else {
				// 如果设置不成功
				callbackResult[p.OutputParam1] = 0
			}
		}
	}
	return callbackResult
}

// -------------------------

// DynamicRegistry 设备的动态注册
func DynamicRegistry() {
	midCredential := conf.Registry.MiddleCredential

	r := register.NewRegister(conf.Registry.ServiceAddress)
	resp, err := r.DynamicRegistry(midCredential)
	if err != nil {
		log.Printf("%s dynamic registry failed, error: %s\n", midCredential, err.Error())
		return
	}
	log.Printf("%s dynamic registry success, ID:%s, device_name:%s, token:%s\n", midCredential, resp.ID, resp.DeviceName, resp.Token)
}

// DynamicRegistryAndConnect 设备的动态注册并上线
func DynamicRegistryAndConnect() {

	options := &mqtt.Options{
		MiddleCredential:       conf.Registry.MiddleCredential,
		DynamocRegisterAddress: conf.Registry.ServiceAddress,
		AutoReconnect:          conf.Device.AutoReconnect,
		LostConnectChan:        make(chan bool),
		ReConnectChan:          make(chan bool),
		Server:                 conf.Mqttbroker.AddressMqtt,
		PropertyType:           constant.PROPERTY_TYPE_BASE,
	}
	m, err := mqtt.InitWithMiddleCredential(options)
	if err != nil {
		log.Panic(err)
	}

	// 连接
	err = m.Connect()
	if err != nil {
		log.Panic(err)
	}

	// 掉线后的处理动作
	go func(o *mqtt.Options) {
		for {
			select {
			case <-options.LostConnectChan:
				// 如果不重连，则退出程序
				if !o.AutoReconnect {
					log.Println("not reconnect to ehub/ihub, procedure will quit!")
					os.Exit(0)
					return
				}
				// 重连，则提示目前暂时掉线
				log.Println("lost connect to ehub/ihub, will auto reconnect!")
			case ok := <-options.ReConnectChan:
				if ok {
					log.Println("设备已连接")
				}
			}
		}
	}(options)
}

func Sub() {
	options := &mqtt.Options{
		Token:           conf.Device.Token,
		AutoReconnect:   conf.Device.AutoReconnect,
		LostConnectChan: make(chan bool),
		ReConnectChan:   make(chan bool),
		Server:          conf.Mqttbroker.AddressMqtt,
		PropertyType:    constant.PROPERTY_TYPE_BASE,
	}

	m, err := mqtt.InitWithToken(options)
	if err != nil {
		log.Panic(err)
	}

	// 连接
	err = m.Connect()
	if err != nil {
		log.Panic(err)
	}

	cb := func(topic string, data []byte) {
		log.Println("=====topic=====:", topic)
		log.Println("=====data=====:", string(data))
	}

	err = m.Subscribe("/test/123456", 0, cb)
	if err != nil {
		log.Panic(err)
	}
}
