# 设备 SDK 文档

### 概述

go 语言 Device SDK 适用于使用 golang 开发业务逻辑的项目。

MQTT 协议版本：3.1.1



### SDK 获取

SDK 1.0

- [最新版本：1.0](https://git.internal.yunify.com/iot-sdk/device-sdk-go)

    

### SDK 安装说明

go get -u -v https://git.internal.yunify.com/iot-sdk/device-sdk-go

Branch: testing

### SDK 使用说明

SDK提供了API供设备厂商调用，用于实现设备属性及事件的上报，以及设备的控制等。

另外，golang语言版本的SDK被设计为可以在不同的操作系统上运行，比如Linux、Windows，因此SDK需要OS或者硬件支持的操作被定义为一些HAL函数，设备厂商在使用SDK开发产品时需要将这些HAL函数进行实现。

产品的业务逻辑、SDK的关系如下图所示：



![image.png](https://upload-images.jianshu.io/upload_images/7998142-8a3090bf9752ef92.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240)



## SDK功能列表

| **模块功能** | **功能点**                                                   |
| :----------: | :----------------------------------------------------------- |
|   设备连云   | 设备可通过该 sdk 与青云IoT物联网平台通信，使用 mqtt 协议进行数据传输，用于设备主动上报信息的场景 |
| 设备身份认证 | token                                                        |
|    物模型    | 使用属性、事件对设备进行描述以及实现，包括：属性上报、事件上报。[物模型](http://103.61.37.229:20080/document/index?document_id=22) |
|   设备控制   | 通过订阅相关 topic，获取下行数据实时控制设备状态             |



### SDK使用示例

1. 设备属性上报(上行)

    ```go
    func TestPubPropertySync(t *testing.T) {
    	options := &index.Options{
    		Token:  "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY3IiOiIxIiwiYXVkIjoiaWFtIiwiYXpwIjoiaWFtIiwiY3VpZCI6ImlhbXItbDZjdGJoZjAiLCJlaXNrIjoiR2VuazUxbm5BLXZyOUJaSnJQQ1gwNnNPSnBabElFZmw4eGlkVEFNbWRjQT0iLCJleHAiOjE2MDU5Njg2MTQsImlhdCI6MTU3NDQzMjYxNCwiaXNzIjoic3RzIiwianRpIjoiVWpJNFdQQW9wNWdQNldPdHJIUTc5USIsIm5iZiI6MCwib3JnaSI6ImlvdGQtZDhjYmEzOTItYWU0NC00MGRmLTk2YzgtNmQ3MWMzMmI4NjZlIiwib3d1ciI6InVzci16eUt5UFNmRyIsInByZWYiOiJxcm46cWluZ2Nsb3VkOmlhbToiLCJydHlwIjoicm9sZSIsInN1YiI6InN0cyIsInRoaWQiOiJpb3R0LVdDcnQ5bk1hUFMiLCJ0eXAiOiJJRCJ9.V0hqewKk6cwwlWzUpBY1HFpMcEvElurmKHh_HtAD816oVsEvl58kK4zpfs1jslASfBLw11OHBE-BD1Zp9FfGicRgTulQ2OUI4t9UiDbmnxGGKODknuP-0lEAb30n6JqLWWZh-rlZlN0tQVixelMC45ftf4LR0OmRH1T250RWO1MNNqqNgral9juTZ8mI9qcvX0yN3Ro7hM_JndeFWc4j9uj_QLus-Sv0mhleMh4i_5uoji7p8XReykwC82Lm2o61EGZZ3T7RCW9GCrSFngIsXnFUxk9mGqUiyW4aqKNkvpcCg-lm3t4fuszc6YW9_YzU53uic14ERRswREf3Wj3vJg",
    		Server: "tcp://127.0.0.1:1883", //192.168.14.120:1883
    	}
    
    	m, err := NewMqtt(options)
    	assert.Nil(t, err)
    	assert.Nil(t, m)
    	time.Sleep(5 * time.Second)
    
    	go func() {
    		var i int64 = 1539362482000
    		for {
    			data := index.Metadata{
    				"CO2Concentration": RandInt64(1, 100),
    				"humidity":         RandInt64(1, 100),
    			}
    
    			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    			reply, err := m.PubPropertySync(ctx, data)
    			cancel()
    			assert.Nil(t, err)
    
    			fmt.Println("PubPropertySync reply", reply)
    
    			i = i + 60000
    			if i > 1570898482000 {
    				i = 1539362482000
    			}
    			time.Sleep(10000 * time.Millisecond)
    		}
    	}()
    	select {}
    }
    ```

    

2. 设备事件上报(上行)

    ```go
    func TestPubEventSync(t *testing.T) {
    	options := &index.Options{
    		Token:  "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY3IiOiIxIiwiYXVkIjoiaWFtIiwiYXpwIjoiaWFtIiwiY3VpZCI6ImlhbXItbDZjdGJoZjAiLCJlaXNrIjoiR2VuazUxbm5BLXZyOUJaSnJQQ1gwNnNPSnBabElFZmw4eGlkVEFNbWRjQT0iLCJleHAiOjE2MDU5Njg2MTQsImlhdCI6MTU3NDQzMjYxNCwiaXNzIjoic3RzIiwianRpIjoiVWpJNFdQQW9wNWdQNldPdHJIUTc5USIsIm5iZiI6MCwib3JnaSI6ImlvdGQtZDhjYmEzOTItYWU0NC00MGRmLTk2YzgtNmQ3MWMzMmI4NjZlIiwib3d1ciI6InVzci16eUt5UFNmRyIsInByZWYiOiJxcm46cWluZ2Nsb3VkOmlhbToiLCJydHlwIjoicm9sZSIsInN1YiI6InN0cyIsInRoaWQiOiJpb3R0LVdDcnQ5bk1hUFMiLCJ0eXAiOiJJRCJ9.V0hqewKk6cwwlWzUpBY1HFpMcEvElurmKHh_HtAD816oVsEvl58kK4zpfs1jslASfBLw11OHBE-BD1Zp9FfGicRgTulQ2OUI4t9UiDbmnxGGKODknuP-0lEAb30n6JqLWWZh-rlZlN0tQVixelMC45ftf4LR0OmRH1T250RWO1MNNqqNgral9juTZ8mI9qcvX0yN3Ro7hM_JndeFWc4j9uj_QLus-Sv0mhleMh4i_5uoji7p8XReykwC82Lm2o61EGZZ3T7RCW9GCrSFngIsXnFUxk9mGqUiyW4aqKNkvpcCg-lm3t4fuszc6YW9_YzU53uic14ERRswREf3Wj3vJg",
    		Server: "tcp://127.0.0.1:1883", // 192.168.14.120:1883
    	}
    	m, err := NewMqtt(options)
    	assert.Nil(t, err)
    	assert.Nil(t, m)
    	time.Sleep(5 * time.Second)
    
    	go func() {
    		var i int64 = 1539362482000
    		for {
    			data := index.Metadata{
    				"int32":  10,
    				"string": "hexing-string",
    				"float":  rand.Float32(),
    				"double": rand.Float64(),
    			}
    			reply, err := m.PubEventSync(context.Background(), "he-event1", data)
    			assert.Nil(t, err)
    			fmt.Println(reply)
    			i = i + 60000
    			if i > 1570898482000 {
    				i = 1539362482000
    			}
    			time.Sleep(10000 * time.Millisecond)
    		}
    	}()
    	select {}
    }
    ```

    



3. 设备控制(下行)

    ```go
    func TestSubDeviceControlSync(t *testing.T) {
    	options := &index.Options{
    		Token:     "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY3IiOiIxIiwiYXVkIjoiaWFtIiwiYXpwIjoiaWFtIiwiY3VpZCI6ImlhbXItaDJjMjl2OXIiLCJlaXNrIjoiVjRSd3NoNjRpcXJhSTVJTHlnZ2xHZFhnV3E1S1JGWWxFYnRwakkxZk9Raz0iLCJleHAiOjE2MDQ1NDc0ODYsImlhdCI6MTU3MzAxMTQ4NiwiaXNzIjoic3RzIiwianRpIjoiVWpJNFdQQW9wNWdQNldPdHJIUTRlQiIsIm5iZiI6MCwib3JnaSI6ImlvdGQtMjk5MDY2MDktNzNiYS00NzBkLWE2ZmQtMGUxYzE3MTkwZmQwIiwib3d1ciI6InVzci1rTFZWQkRxZCIsInByZWYiOiJxcm46cWluZ2Nsb3VkOmlhbToiLCJydHlwIjoicm9sZSIsInN1YiI6InN0cyIsInRoaWQiOiJpb3R0LWZUeXRjS1BWTlEiLCJ0eXAiOiJJRCJ9.C1oCwaviLAsmb42mDXX4mXw2h0ccXYV8Kd5mAGkCpxpOFM7Rd7lOL2kGMJpvv_I5caOTlSNiFwMe2L2eXiA_dsZPBEW08dmzghLZXpVABFG7KJOrxT5t6WBYzVCOezt4CynSXheIs0NjSMZ5VBTdiEjj8GIi5iAIWUaYrEeFOlj3IZPp7ddr82rkog9OIDnHDvyXDK2MruKAb7xZ2QZFa0Wg1GKixFUhfT0iU37pQZbsGAduj-kB9z4o_ZwtP8gFko6AkW8WuBzzXhs35cQty2vXJ3ohxKnXtoiwChNfIQmNr8Cc7VJmQTmrQPrgmK3uMnxi02SQXsF2vd0HmpA_7A",
    		Server:    "tcp://127.0.0.1:1883", // tcp://192.168.14.120:1883 tcp://192.168.14.120:8055
    		Identifer: "start",
    	}
    	m, err := NewMqtt(options)
    	assert.Nil(t, err)
    
    	go m.SubDeviceControlSync()
    	fmt.Println("run")
    
    	time.Sleep(15 * time.Second)
    
    	m.UnSubDeviceControlSync()
    
    	time.Sleep(3 * time.Second)
    }
    ```

    



## 历史版本清单

| **版本号** | **发布日期** | **下载链接** | **更新内容**                                                 |
| :--------- | :----------- | :----------- | :----------------------------------------------------------- |
| 1.0        | 2020/02/07   |              | 读取设备凭证：手动拷贝到设备上，替换示例程序中的变量；<br />端设备连接、收发消息消息、重连<br />边设备连接、收发消息消息、重连<br /> |