# 设备 SDK 文档

### 概述

----------------------

这篇文档介绍了如何安装和配置 设备sdk，以及提供了相关例子来演示如何使用 设备sdk 上报设备数据以及控制设备；

支持MQTT 协议版本：3.1.1



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

3. 了解青云的物模型

    http://103.61.37.229:20080/document/index?document_id=22

4. 创建模型

    TODO:

    青云平台创建模型，获得属性名称

5. 注册设备，绑定模型

    TODO:

    青云平台，获取设备凭证(token)，将上面创建的物模型和设备进行绑定

#### 示例1：设备直连 ihub(青云iot平台)

1. 设备连接

    ```go
    // 设备连接，token 为设备凭证，Server 为青云 iot 平台 ihub
    func main() {
    	options := &index.Options{
    		Token:  "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY3IiOiIxIiwiYXVkIjoiaWFtIiwiYXpwIjoiaWFtIiwiY3VpZCI6ImlhbXItdW9nbjh3N3IiLCJlaXNrIjoiSi1hSHdySDgwTHVXWEI1bGJMQ2E3cm1uTFVfZ1RlbWhJd0FJdVN6T29qZz0iLCJleHAiOjE2MTI5MTk3MjAsImlhdCI6MTU4MTM4MzcyMCwiaXNzIjoic3RzIiwianRpIjoiVWpJNFdQQW9wNWdQNldPdHJIUU5rUyIsIm5iZiI6MCwib3JnaSI6ImlvdGQtYjViODJiYzMtZWEwYS00NDRjLTliYjMtNDkzOWE0YTgzMmRiIiwib3d1ciI6InVzci1rZUF5dG16MSIsInByZWYiOiJxcm46cWluZ2Nsb3VkOmlhbToiLCJydHlwIjoicm9sZSIsInN1YiI6InN0cyIsInRoaWQiOiJpb3R0LXlWQXd4OXJiOGoiLCJ0eXAiOiJJRCJ9.Hq6zIwQCpBV897fsVl-WHxBtgtH8xe8umcp5QIQ3p1lSHrYUV_ofrbJ5oZKasUKwYqxlhhcxjX2f3U9OCCvOj8yGjIyK8vrf8vJBbNwW48fkCiVnFOpoKui8k9Fg13qNl0AUD8TmOWAukn3uQTI7gKW6fhwmWkdZD8cLOraEBvkGkrL19Nlw-JuU-MWXeNB2p1F5CahUAvDD78zUHkPJTZ-X3v9d73YyUlSV2CrJvAJGpae6sHCXk4iS8KyPQ1GNjPmjD9qBdbzr5cdIA3LjIkuppaWb0i8vymhvLaqcfD5EnEfu8aKNNLGBedEI3c8BlXOLSgp5_BldOJuP2GnGfQ",
    		Server: "tcp://192.168.14.120:8055",
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
    ```

    设备连接后可在青云 iot 平台查看设备的连接状态

2. 属性上报

    ```go
    func main() {
    	options := &index.Options{
    		Token:  "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY3IiOiIxIiwiYXVkIjoiaWFtIiwiYXpwIjoiaWFtIiwiY3VpZCI6ImlhbXItdW9nbjh3N3IiLCJlaXNrIjoiSi1hSHdySDgwTHVXWEI1bGJMQ2E3cm1uTFVfZ1RlbWhJd0FJdVN6T29qZz0iLCJleHAiOjE2MTI5MTk3MjAsImlhdCI6MTU4MTM4MzcyMCwiaXNzIjoic3RzIiwianRpIjoiVWpJNFdQQW9wNWdQNldPdHJIUU5rUyIsIm5iZiI6MCwib3JnaSI6ImlvdGQtYjViODJiYzMtZWEwYS00NDRjLTliYjMtNDkzOWE0YTgzMmRiIiwib3d1ciI6InVzci1rZUF5dG16MSIsInByZWYiOiJxcm46cWluZ2Nsb3VkOmlhbToiLCJydHlwIjoicm9sZSIsInN1YiI6InN0cyIsInRoaWQiOiJpb3R0LXlWQXd4OXJiOGoiLCJ0eXAiOiJJRCJ9.Hq6zIwQCpBV897fsVl-WHxBtgtH8xe8umcp5QIQ3p1lSHrYUV_ofrbJ5oZKasUKwYqxlhhcxjX2f3U9OCCvOj8yGjIyK8vrf8vJBbNwW48fkCiVnFOpoKui8k9Fg13qNl0AUD8TmOWAukn3uQTI7gKW6fhwmWkdZD8cLOraEBvkGkrL19Nlw-JuU-MWXeNB2p1F5CahUAvDD78zUHkPJTZ-X3v9d73YyUlSV2CrJvAJGpae6sHCXk4iS8KyPQ1GNjPmjD9qBdbzr5cdIA3LjIkuppaWb0i8vymhvLaqcfD5EnEfu8aKNNLGBedEI3c8BlXOLSgp5_BldOJuP2GnGfQ",
    		Server: "tcp://192.168.14.120:8055", // 192.168.14.120:8055
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
    
        fmt.Printf("PubProperty reply:%+v", reply)
    }
    ```

    上述代码中，CO2Concentration、humidity 分别是 **前置条件3** 中创建的模型属性；

    属性上报成功后，将会在 青云iot 平台 设备界面显示上报数据信息；

3. 事件上报

    ```go
    func main() {
    	options := &index.Options{
    		Token:  "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY3IiOiIxIiwiYXVkIjoiaWFtIiwiYXpwIjoiaWFtIiwiY3VpZCI6ImlhbXItdW9nbjh3N3IiLCJlaXNrIjoiSi1hSHdySDgwTHVXWEI1bGJMQ2E3cm1uTFVfZ1RlbWhJd0FJdVN6T29qZz0iLCJleHAiOjE2MTI5MTk3MjAsImlhdCI6MTU4MTM4MzcyMCwiaXNzIjoic3RzIiwianRpIjoiVWpJNFdQQW9wNWdQNldPdHJIUU5rUyIsIm5iZiI6MCwib3JnaSI6ImlvdGQtYjViODJiYzMtZWEwYS00NDRjLTliYjMtNDkzOWE0YTgzMmRiIiwib3d1ciI6InVzci1rZUF5dG16MSIsInByZWYiOiJxcm46cWluZ2Nsb3VkOmlhbToiLCJydHlwIjoicm9sZSIsInN1YiI6InN0cyIsInRoaWQiOiJpb3R0LXlWQXd4OXJiOGoiLCJ0eXAiOiJJRCJ9.Hq6zIwQCpBV897fsVl-WHxBtgtH8xe8umcp5QIQ3p1lSHrYUV_ofrbJ5oZKasUKwYqxlhhcxjX2f3U9OCCvOj8yGjIyK8vrf8vJBbNwW48fkCiVnFOpoKui8k9Fg13qNl0AUD8TmOWAukn3uQTI7gKW6fhwmWkdZD8cLOraEBvkGkrL19Nlw-JuU-MWXeNB2p1F5CahUAvDD78zUHkPJTZ-X3v9d73YyUlSV2CrJvAJGpae6sHCXk4iS8KyPQ1GNjPmjD9qBdbzr5cdIA3LjIkuppaWb0i8vymhvLaqcfD5EnEfu8aKNNLGBedEI3c8BlXOLSgp5_BldOJuP2GnGfQ",
    		Server: "tcp://192.168.14.120:8055", // 192.168.14.120:8055
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
    		"string": "hexing-string",
    		"float":  rand.Float32(),
    		"double": rand.Float64(),
    	}
    	reply, err := m.PubEvent(context.Background(), "he-event1", data)
    	if err != nil {
    		panic(err)
    	}
    	fmt.Printf("PubEvent reply:%+v\n", reply)
    }
    ```

    

4. 设备控制(下行)

    ```go
    func main() {
    	options := &index.Options{
    		Token:  "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY3IiOiIxIiwiYXVkIjoiaWFtIiwiYXpwIjoiaWFtIiwiY3VpZCI6ImlhbXItdW9nbjh3N3IiLCJlaXNrIjoiSi1hSHdySDgwTHVXWEI1bGJMQ2E3cm1uTFVfZ1RlbWhJd0FJdVN6T29qZz0iLCJleHAiOjE2MTI5MTk3MjAsImlhdCI6MTU4MTM4MzcyMCwiaXNzIjoic3RzIiwianRpIjoiVWpJNFdQQW9wNWdQNldPdHJIUU5rUyIsIm5iZiI6MCwib3JnaSI6ImlvdGQtYjViODJiYzMtZWEwYS00NDRjLTliYjMtNDkzOWE0YTgzMmRiIiwib3d1ciI6InVzci1rZUF5dG16MSIsInByZWYiOiJxcm46cWluZ2Nsb3VkOmlhbToiLCJydHlwIjoicm9sZSIsInN1YiI6InN0cyIsInRoaWQiOiJpb3R0LXlWQXd4OXJiOGoiLCJ0eXAiOiJJRCJ9.Hq6zIwQCpBV897fsVl-WHxBtgtH8xe8umcp5QIQ3p1lSHrYUV_ofrbJ5oZKasUKwYqxlhhcxjX2f3U9OCCvOj8yGjIyK8vrf8vJBbNwW48fkCiVnFOpoKui8k9Fg13qNl0AUD8TmOWAukn3uQTI7gKW6fhwmWkdZD8cLOraEBvkGkrL19Nlw-JuU-MWXeNB2p1F5CahUAvDD78zUHkPJTZ-X3v9d73YyUlSV2CrJvAJGpae6sHCXk4iS8KyPQ1GNjPmjD9qBdbzr5cdIA3LjIkuppaWb0i8vymhvLaqcfD5EnEfu8aKNNLGBedEI3c8BlXOLSgp5_BldOJuP2GnGfQ",
    		Server: "tcp://192.168.14.120:8055",
            Identifer: "start",
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
    ```

    Identifer：设备控制标识符

    运行上述代码，在 mqttbox 中 **publish** 打印出来的 topic，在 payload 中输入数据模型：

    ```go
    {
        "id": "123",
        "version": "1.0",
        "params": {
            "label":"on",
            "image":"23.6"
        }
    }
    ```

    **subscribe** 上述topic + "\_reply"，会获取如下数据结构：

    ```go
    {
        "id": "123",
        "code": 200,
        "data": {
            "label":"on",
            "image":"23.6"
        }
    }
    ```



## 历史版本清单

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