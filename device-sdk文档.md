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
        
        

#### 示例1：设备直连 ihub(青云iot平台)

##### 准备工作

1. 了解青云的物模型

    http://103.61.37.229:20080/document/index?document_id=22

2. 在青云iot平台创建模型

    TODO:

    青云平台创建模型，获得属性名称

3. 在青云iot平台注册设备，绑定模型

    TODO:

    青云平台，获取设备凭证(token)，将上面创建的物模型和设备进行绑定

#####  代码演示

1. /设备连接

    ```go
    // 设备连接，token 为设备凭证，Server 为青云 iot 平台 ihub
    func main() {
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
    ```
    设备连接后可在青云 iot 平台查看设备的连接状态

    ![image.png](https://upload-images.jianshu.io/upload_images/7998142-25f110a214a9b64e.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240)

2. 属性上报

    ```go
    func main() {
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
    ```

    上述代码中，CO2Concentration、humidity 分别是 **前置条件3** 中创建的模型属性；
    属性上报成功后，将会在 青云iot 平台 设备界面的属性模块显示上报数据信息：
    
    ![image.png](https://upload-images.jianshu.io/upload_images/7998142-944efe481c009482.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240)


3. 事件上报

    ```go
    func main() {
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
    ```
    PubEvent 函数中第三个参数为 事件的 identifier；
    事件上报成功后，将会在 青云iot 平台 设备界面的事件模块显示上报数据信息：

    ![image.png](https://upload-images.jianshu.io/upload_images/7998142-a3834da7f90d116d.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240)

    
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