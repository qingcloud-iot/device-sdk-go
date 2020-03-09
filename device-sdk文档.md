# 设备 SDK 文档

### 概述

----------------------

这篇文档介绍了如何安装和配置 设备sdk，以及提供了相关例子来演示如何使用 设备sdk 上报设备数据以及控制设备；

支持MQTT 协议版本：3.1.1
golang version：1.13及以上



### SDK 获取

-------

- [最新版本：1.0](https://git.internal.yunify.com/iot-sdk/device-sdk-go)



### SDK 功能列表

-------------

| **模块功能** | **功能点**                                                   |
| :----------: | :----------------------------------------------------------- |
|   设备连云   | 设备可通过该 sdk 与青云IoT物联网平台通信，使用 mqtt 协议进行数据传输，用于设备主动上报信息的场景 |
| 设备身份认证 | token(即设备凭证)                                            |
|   属性上报   | 向特定 topic 上报设备属性数据 [物模型](http://103.61.37.229:20080/document/index?document_id=22) |
|   事件上报   | 向特定 topic上报设备事件                                     |
|   设备控制   | 通过订阅相关 topic，获取下行数据实时控制设备状态             |



### SDK使用示例

------------------------

#### 前置条件

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

        vim ~/.zshrc

        ```go
        export GOPATH=$HOME/gopath
        export PATH=$PATH:/usr/local/go/bin:$GOPATH/bin
        ```

        source ~/.zshrc
        
        

#### 示例1：设备直连 ihub(青云iot平台)

##### 准备工作

1. 环境配置
    - /etc/hosts 添加：
        192.168.14.179     console-staging.qingcloud.com iot-staging.qingcloud.com
        192.168.14.121     iot-api.qingcloud.com
    - 连接公司 vpn
    - 地址及账户
        http://iot-staging.qingcloud.com
        rainsong@yunify.com
        zhu88jie

2. 了解青云的物模型

    http://103.61.37.229:20080/document/index?document_id=22

3. 在青云iot平台创建模型
    本例中使用物模型：EdgeWise

    ```go
    {
        "id": "iott-yVAwx9rb8j",
        "name": "EdgeWize",
        "type": 2,
        "icon": "cpe",
        "user_id": "usr-keAytmz1",
        "description": "云锡：边缘模型",
        "properties": {
            "AlarmState": {
                "property_id": "iotp-c7c90665-2d62-483b-aa7c-af5f2d1fbe93",
                "property_name": "",
                "identifier": "AlarmState",
                "access": "all",
                "type": "bool",
                "define": {},
                "description": "告警状态"
            },
            "MaxValue": {
                "property_id": "iotp-14337e33-4772-410b-a3e0-55950ed758ee",
                "property_name": "",
                "identifier": "MaxValue",
                "access": "all",
                "type": "float",
                "define": {},
                "description": "最大值"
            }
        },
        "events": {
            "statistics": {
                "event_id": "iote-790b0ea1-ac2c-4179-85f2-9c76f9982fc8",
                "identifier": "statistics",
                "event_name": "数值统计",
                "description": "最大最小值统计",
                "type": "info",
                "output": [
                    {
                        "id": "a718a576-40ce-44b7-841d-30313b5cdbd5",
                        "identifier": "max",
                        "output_name": "最大值",
                        "type": "float",
                        "define": {}
                    },
                    {
                        "id": "76f1ec21-2793-481c-87b0-3e94a21f33dd",
                        "identifier": "min",
                        "output_name": "最小值",
                        "type": "float",
                        "define": {}
                    }
                ]
            },
            "threshold": {
                "event_id": "iote-f4e028dc-b6fc-437c-9bec-9a947c52f79b",
                "identifier": "threshold",
                "event_name": "阈值告警",
                "description": "阈值告警事件",
                "type": "info",
                "output": [
                    {
                        "id": "c61c54b2-c3c4-460f-8eac-feabc24c60e1",
                        "identifier": "identifier",
                        "output_name": "告警点位",
                        "type": "string",
                        "define": {}
                    },
                    {
                        "id": "e36bb956-dc18-447f-8ae9-f99312c27b3d",
                        "identifier": "val",
                        "output_name": "告警值",
                        "type": "float",
                        "define": {}
                    }
                ]
            }
        },
        "actions": {
            "connect": {
                "action_id": "iots-23ff91e2-d8a3-4e8d-87d7-acd25c3e2dc6",
                "identifier": "connect",
                "action_name": "连接opc",
                "description": "连接opc server",
                "call_type": "sync",
                "input": [],
                "output": [
                    {
                        "id": "18ae0993-0b52-42af-b96b-b64b7f40cc65",
                        "identifier": "reply",
                        "name": "",
                        "type": "string",
                        "define": {}
                    }
                ]
            },
            "disconnect": {
                "action_id": "iots-a0ad97a0-c16d-4ec1-a276-c5c0b9a95930",
                "identifier": "disconnect",
                "action_name": "断开opc连接",
                "description": "断开opc连接",
                "call_type": "sync",
                "input": [],
                "output": [
                    {
                        "id": "447ca1f6-472c-4f21-9e58-2777d112749e",
                        "identifier": "reply",
                        "name": "",
                        "type": "string",
                        "define": {}
                    }
                ]
            }
        }
    }
    ```

    青云平台创建模型，可以得到属性名称、属性类型、事件identifier、控制identifier等；

4. 在青云iot平台注册设备，绑定模型
    本例使用设备：TestEdge
    在青云平台注册设备，以及将上面创建的物模型和设备进行绑定，获取设备凭证(token)；

#####  代码演示

1. 设备连接

    ```go
    package main

    import (
        "context"
        "fmt"
        "time"

        "git.internal.yunify.com/iot-sdk/device-sdk-go/index"
        "git.internal.yunify.com/iot-sdk/device-sdk-go/mqtt"
    )

    // 设备连接，token 为设备凭证，Server 为青云 iot 平台 ihub
    func main() {
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
    ```
    运行上述代码前后，即可在青云 iot 平台查看设备连接状态，如：

    ![image.png](https://upload-images.jianshu.io/upload_images/7998142-dd7c5790db12afe5.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240)

    ![image.png](https://upload-images.jianshu.io/upload_images/7998142-25f110a214a9b64e.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240)

2. 属性上报
    用于将设备相关数据上报，可在青云iot 平台实时查看设备运行动态；
    ```go
    package main

    import (
        "context"
        "fmt"
        "time"

        "git.internal.yunify.com/iot-sdk/device-sdk-go/index"
        "git.internal.yunify.com/iot-sdk/device-sdk-go/mqtt"
    )

    func main() {
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
    ```

    上述代码中，CO2Concentration、humidity 分别是 **前置条件3** 中创建的模型属性；
    属性上报成功后，将会在 青云iot 平台 设备界面的属性模块显示上报数据信息：
    
    ![image.png](https://upload-images.jianshu.io/upload_images/7998142-944efe481c009482.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240)


3. 事件上报

    ```go
    package main

    import (
        "context"
        "fmt"
        "time"

        "git.internal.yunify.com/iot-sdk/device-sdk-go/index"
        "git.internal.yunify.com/iot-sdk/device-sdk-go/mqtt"
    )

    func main() {
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
    ```
    PubEvent 函数中第三个参数为 事件的 identifier；
    事件上报成功后，将会在 青云iot 平台 设备界面的事件模块显示上报数据信息：

    ![image.png](https://upload-images.jianshu.io/upload_images/7998142-a3834da7f90d116d.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240)

4. 设备控制(下行)

    ```go
    package main

    import (
        "context"
        "fmt"
        "time"

        "git.internal.yunify.com/iot-sdk/device-sdk-go/index"
        "git.internal.yunify.com/iot-sdk/device-sdk-go/mqtt"
    )

    func main() {
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
    ```

    Identifer：设备控制标识符;
    通过接口下发数据，并可以得到响应结果：
    Identifer：设备标识符
    id：设备id
    Params：根据 准备工作1 中的青云物模型的设备控制模型 和 准备工作2 中新建的控制的数据模可以得到

    ![image.png](https://upload-images.jianshu.io/upload_images/7998142-079de02295e411a6.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240)

    ![image.png](https://upload-images.jianshu.io/upload_images/7998142-b1d225dc7c3786d1.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240)

    或使用命令行命令 curl：
    ```go
    curl -X POST "http://iot-api.qingcloud.com:8889/api/v1/devices/iotd-59f685ce-70f9-4485-9985-271d5dfb9475/call/connect" -H "accept: application/json" -H "Authorization: QCUUCDRNSJECJRCMOMPPHP:signaturea" -H "Content-Type: application/json" -d "{ \"params\": { \"reply\":\"this is control test\" }, \"thing_id\": \"iott-yVAwx9rb8j\"}"

    reply: {"code":"0","data":{"reply":"this is control test"}}
    ```
    或运行：
    ```go
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
    ```
    便可在程序端收到下行消息


### 历史版本清单
-------------
| **版本号** | **发布日期** | **下载链接** | **更新内容**                                                 |
| :--------- | :----------- | :----------- | :----------------------------------------------------------- |
| 1.0        | 2020/02/07   |              | 读取设备凭证：手动拷贝到设备上，替换示例程序中的变量；<br />端设备连接、收发消息消息、重连<br />边设备连接、收发消息消息、重连<br /> |



### 附录

-------

#### 1. SDK 的本地辅助测试

   - mosquitto 或 ehub

     - mosquitto

         sudo apt-get install mosquitto

         sudo service mosquitto start 

         sudo service mosquitto stop

         sudo service mosquitto status

     - ehub

         https://git.internal.yunify.com/edge/exia

- mqttbox

    下载：http://workswithweb.com/html/mqttbox/installing_apps.html

    使用：[MQTT系列教程3（客户端工具MQTTBox的安装和使用）](https://www.hangge.com/blog/cache/detail_2350.html)

-----