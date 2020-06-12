本章节介绍了如何安装和配置**设备 SDK**，以及提供了相关例子来演示如何使用**设备 SDK**上报设备数据以及服务调用；

支持MQTT 协议版本：3.1.1

golang version：1.13及以上  

### SDK 获取

-------

- 最新版本：[1.0](https://iot-sdk.pek3a.qingstor.com/v1.0/device-sdk-go.tar.gz)

### SDK 功能列表

-------------

|    **模块功能**    | **功能点**                                                   |
| :----------------: | :----------------------------------------------------------- |
|      设备连云      | 设备可通过该 SDK 与青云IoT物联网平台通信，使用 mqtt/mqtts 协议进行数据传输，用于设备主动上报信息的场景 |
|      自动重连      | 由于网络等原因，设备存在掉线的可能，用户可以通过配置控制设备掉线后是否需要重新连接云平台 |
|    设备身份认证    | token(设备凭证)                                              |
|      属性上报      | 向特定 topic 上报设备属性数据                                |
|      事件上报      | 向特定 topic上报设备事件                                     |
|      设备控制      | 通过订阅相关 topic，获取下行数据实时控制设备状态             |
|      动态注册      | 利用中间凭证实现大批量设备接入青云 IoT 物联网平台            |
| 动态注册并设备连云 | 使用中间凭证动态注册后，通过获得的设备 token 直接连云        |

### SDK API 列表

-------------

|         函数         | 功能                                       |
| :------------------: | :----------------------------------------- |
|       Options        | 初始化 SDK Client 相关选项                 |
| DeviceControlHandler | 服务调用结构体，用于处理下行数据的业务逻辑 |
|       Connect        | 设备连接物联网平台                         |
|      DisConnect      | 设备取消连接物联网平台                     |
|     PubProperty      | 推送设备属性                               |
|       PubEvent       | 推送设备事件                               |
|   SubDeviceControl   | 订阅 topic   |
|  UnSubDeviceControl  | 取消订阅                                   |



### SDK使用简介

------------------------

#### 1. 前置条件

1. 系统

    linux/win/macos，1核2G即可；

2. go 环境

    - 下载安装包  

        wget https://dl.google.com/go/go1.13.linux-amd64.tar.gz

    - 解压并移动到指定目录   

        tar -xvzf go1.13.linux-amd64.tar.gz  

        sudo mv go /usr/local/  

    - 建立 go 的工作空间  

        在/home目录下, 建立一个 gopath目录，然后建立三个子目录：src、pkg、bin  

        src — 里面每一个子目录，就是一个包。包内是Go的源码文件  

        pkg — 编译后生成的，包的目标文件  

        bin — 生成的可执行文件

    - 设置 GOPATH 环境变量  

        vim ~/.bashrc

        ```go
        export GOPATH=$HOME/gopath
        export PATH=$PATH:/usr/local/go/bin:$GOPATH/bin
        ```

        source ~/.bashrc  

#### 2. 配置文件

config.yml

```yaml
device:
    token: <your_device_token>
    auto_reconnect: true
mqttbroker:
    address_mqtt: <your_gateway_address>
    address_mqtts: <your_tls_gateway_address>
registry:
    middle_credential: <your_middle_credential>
    service_address: <your_service_address>
```

- device

    token: 设备凭证，注册设备时可获取到，解析 token 可获取到 设备ID、数据模型ID 等信息；

    auto_reconnect: 设置设备掉线后是否自动重连，true / false；

- mqttbroker

    设备数据上报的目的地址，可以是边端，也可以是云端；

    address_mqtt: 通过 mqtt 方式上报数据

    address_mqtts: 通过 mqtts 方式加密上报数据

- registry.

    middle_credential: 中间凭证，大批量设备注册后会产生一个中间凭证，通过该凭证能够实现同批次的海量设备使用相同的方式和信息接入平台。同时还可以在动态注册之后获得专属的设备凭证，也就是上面的 token；

    service_address: 动态注册的服务地址

#### 3. 设备接入

- 通过 token 接入

    options 中传入参数为配置文件中的 token、mqttbroker.address

    ```go
    options := &mqtt.Options{
        Token:     conf.Device.Token,
        AutoReconnect:   conf.Device.AutoReconnect,
		LostConnectChan: make(chan bool),
        Server:    conf.Mqttbroker.Address,

        KeepAlive:            60,               // 心跳间隔, 默认 30s
	    MaxReconnectInterval: 20 * time.Minute, // 最大重连间隔, 默认 10 * time.Minute
    }
    
    // 初始化
    m, err := mqtt.InitWithToken(options)
    if err != nil {
        panic(err)
    }

    // 连接
    err = m.Connect()
    if err != nil {
        panic(err)
    }

    // 掉线(离线)后的处理动作
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
    ```

    设备接入云端或边端后，才能上报数据和服务调用；

- 通过中间凭证接入

    除了通过设备专属凭证接入外，还可以通过批量设备的中间凭证接入

    options 中传入参数为配置文件中的 middle_credential、service_address

    ```go
    options := &mqtt.Options{
        MiddleCredential:       conf.Registry.MiddleCredential,
        DynamocRegisterAddress: conf.Registry.ServiceAddress,
        AutoReconnect:   conf.Device.AutoReconnect,
        LostConnectChan: make(chan bool),
        Server:       conf.Mqttbroker.Address,
        PropertyType: constant.PROPERTY_TYPE_BASE,
    }

    m, err := mqtt.InitWithMiddleCredential(options)
    if err != nil {
        panic(err)
    }
    
    // 连接
    err = client.Connect()
    if err != nil {
    panic(err)
    }

    // 掉线(离线)后的处理动作
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
    ```

如果要使用 mqtts 进行加密通信，将 options 的 CertFilePath 字段设置为证书地址即可！

[设备接入使用示例](https://iot-docs.qingcloud.com/beta/zh-CN/tutorials/access-end-devices/)

#### 4. 属性上报

通过 PubProperty 方法上报设备属性，传入参数 propertyData 为**数据模型**中定义的属性 identifier 及属性值

```
propertyData := define.PropertyKV{
		"temp": 11,
	}
ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
reply, err := client.PubProperty(ctx, propertyData)
```

属性上报后，会响应 reply 给用户，如果 reply.Code 等于 200，则表示上报成功，否则上报失败，失败信息可在 reply.Data 中查看；

属性上报成功后，可以在物联网平台查看上报的属性值；

[属性上报使用示例](https://iot-docs.qingcloud.com/beta/zh-CN/tutorials/send-data/#%E4%B8%8A%E6%8A%A5%E5%B1%9E%E6%80%A7%E6%95%B0%E6%8D%AE)

#### 5. 事件上报

通过 PubEvent 方法上报事件，传入参数 PubEvent 为**数据模型**中定义的事件 identifier 及事件信息

```
eventData := define.PropertyKV{
    "temperature": float64(DeviceTemprature),
    "reason":      reason,
}
ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
reply, err := client.PubEvent(ctx, eventData, eventIdentifier)
```

属性上报后，会响应 reply 给用户，如果 reply.Code 等于 200，则表示上报成功，否则上报失败，失败信息可在 reply.Data 中查看；

事件上报成功后，可以在物联网平台查看上报的事件信息；

[事件上报使用示例](https://iot-docs.qingcloud.com/beta/zh-CN/tutorials/send-data/#%E4%B8%8A%E6%8A%A5%E4%BA%8B%E4%BB%B6%E6%95%B0%E6%8D%AE)

#### 6. 服务调用

服务调用即服务端对设备属性值进行设置，用户需要对服务端下发的数据根据需求进行处理；

```go
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

type InAndOutputParameters struct {
	InputParam1  string
	OutputParam1 string
	OutputParam2 string
}

// DeviceControlCallback 服务调用的回调函数, 实现 handler 方法即可
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
```

[服务调用使用示例](https://iot-docs.qingcloud.com/beta/zh-CN/tutorials/send-data/#%E8%B0%83%E7%94%A8%E6%9C%8D%E5%8A%A1)

#### 7. 动态注册

通过 DynamicRegistry 方法进行批量设备的动态注册

```go
midCredential := conf.Registry.MiddleCredential

r := register.NewRegister(conf.Registry.ServiceAddress)
resp, err := r.DynamicRegistry(midCredential)
if err != nil {
	fmt.Printf("%s dynamic registry failed, error: %s\n", midCredential, err.Error())
	return
}
```

注册成功后，可以通过 resp.Name 获取设备名，通过 resp.Token 获取设备的专属凭证；

[动态注册使用示例](https://iot-docs.qingcloud.com/beta/zh-CN/use-guide/dev-token/#%E5%8A%A8%E6%80%81%E6%B3%A8%E5%86%8C%E8%AE%BE%E5%A4%87)


### 历史版本清单

-------------

| **版本号** | **发布日期** | **下载链接** | **更新内容**                                                 |
| :--------- | :----------- | :----------- | :----------------------------------------------------------- |
| 1.0        | 2020/02/07   |              | 读取设备凭证：手动拷贝到设备上，替换示例程序中的变量；<br />端设备连接、收发消息消息、重连<br />边设备连接、收发消息消息、重连<br /> |



### 附录

-------

#### 1. SDK 的本地辅助测试

- mosquitto

    - mosquitto

        sudo apt-get install mosquitto

        sudo service mosquitto start 

        sudo service mosquitto stop

        sudo service mosquitto status

- mqttbox

    下载：http://workswithweb.com/html/mqttbox/installing_apps.html

    使用：[MQTT系列教程3（客户端工具MQTTBox的安装和使用）](https://www.hangge.com/blog/cache/detail_2350.html)