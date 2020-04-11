package mqtt

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"git.internal.yunify.com/iot-sdk/device-sdk-go/internal/register"
	"time"

	iClient "git.internal.yunify.com/iot-sdk/device-sdk-go/internal/client"
	"git.internal.yunify.com/iot-sdk/device-sdk-go/internal/constant"
	"git.internal.yunify.com/iot-sdk/device-sdk-go/internal/define"

	mqttp "github.com/eclipse/paho.mqtt.golang"
)

//
type Handler func(inputIdentifier string, msg *define.Message) error

// DeviceControlHandler 服务调用结构体，用于处理下行数据的业务逻辑
type DeviceControlHandler struct {
	ServiceIdentifer string
	InputIdentifier  string
	ServiceHandler   Handler
}

type Options struct {
	Token string // 权限验证，及获取 ModelID、EntityID

	MiddleCredential       string // 批量设备注册的中间凭证
	DynamocRegisterAddress string // 动态注册的服务地址

	Server         string                 // mqtt server
	PropertyType   string                 // 属性分组（系统属性platform、基础属性base）
	DeviceHandlers []DeviceControlHandler // 服务调用的回调函数
	EntityID       string                 // 设备 id
	ModelID        string                 // 模型 id
}

type MqttClient struct {
	Client mqttp.Client

	EntityId string
	ModelId  string

	PropertyType    string
	UnSubScribeChan chan bool
}

func initMQTTClient(options *Options) mqttp.Client {

	opts := mqttp.NewClientOptions()
	opts.AddBroker("tcp://" + options.Server)
	if options.EntityID != "" {
		opts.SetClientID(options.EntityID)
		opts.SetUsername(options.EntityID)
	}
	opts.SetPassword(options.Token)
	opts.SetCleanSession(true)
	opts.SetAutoReconnect(true)
	opts.SetKeepAlive(30 * time.Second)
	opts.SetConnectionLostHandler(func(client mqttp.Client, e error) {
		fmt.Println("lost connect!")
	})
	opts.SetOnConnectHandler(func(client mqttp.Client) {
		fmt.Println("connect ehub/ihub success!")
	})

	if options.DeviceHandlers != nil {
		opts.SetDefaultPublishHandler(func(client mqttp.Client, msg mqttp.Message) {
			fmt.Printf("[sdk-go sub] topic: %s, paload: %s\n", msg.Topic(), string(msg.Payload()))

			for _, handler := range options.DeviceHandlers {
				switch {
				case msg.Topic() == BuildServiceControlReply(options.ModelID, options.EntityID, handler.ServiceIdentifer):
					var err error
					topic := msg.Topic()
					payload := msg.Payload()

					message, err := ParseMessage(payload)
					if err != nil {
						return
					}

					// 执行回调函数进行服务调用
					if err = handler.ServiceHandler(handler.InputIdentifier, message); err != nil {
						fmt.Printf("topic:%s, execute callback error: %s\n", topic, err.Error())
					}

					// reply
					if err = Reply(message, client, topic); err != nil {
						fmt.Printf("topic:%s, reply error: %s\n", topic, err.Error())
						return
					}
				default:
				}
			}
		})
	}

	mqttClient := mqttp.NewClient(opts)
	return mqttClient
}

// InitWithToken 使用 token 进行设备通讯
func InitWithToken(options *Options) (iClient.Client, error) {

	m := &MqttClient{
		UnSubScribeChan: make(chan bool),
	}

	if options.Token == "" {
		return nil, fmt.Errorf("token can not be blank")
	}

	entityID, modelID, err := ParseToken(options.Token)
	if err != nil {
		return nil, fmt.Errorf("token is invalid: %s", err.Error())
	}
	options.EntityID = entityID
	options.ModelID = modelID
	mqttClient := initMQTTClient(options)

	m.Client = mqttClient
	m.EntityId = entityID
	m.ModelId = modelID

	if options.PropertyType != "" {
		m.PropertyType = options.PropertyType
	}
	return m, nil
}

// InitWithMiddleCredential 使用中间凭证进行设备通讯
func InitWithMiddleCredential(options *Options) (iClient.Client, error) {

	if options.MiddleCredential == "" {
		return nil, fmt.Errorf("MiddleCredential can not be blank")
	}

	if options.DynamocRegisterAddress == "" {
		return nil, fmt.Errorf("DynamocRegisterAddress can not be blank")
	}

	// 通过 middleCredential 进行动态注册，获取设备 token
	r := register.NewRegister(options.DynamocRegisterAddress)
	resp, err := r.DynamicRegistry(options.MiddleCredential)
	if err != nil {
		return nil, err
	}

	options.Token = resp.Token
	return InitWithToken(options)
}

// Connect 连接 ihub 或 ehub
func (m *MqttClient) Connect() error {
	if token := m.Client.Connect(); !token.WaitTimeout(5*time.Second) || token.Error() != nil {
		return fmt.Errorf("mqtt client connect fail")
	}
	return nil
}

// DisConnect 断开连接 ihub 或 ehub
func (m *MqttClient) DisConnect() {
	m.Client.Disconnect(QUIESCE)
}

// PubProperty 上报属性
func (m *MqttClient) PubProperty(ctx context.Context, meta define.PropertyKV) (*define.Reply, error) {
	reply := &define.Reply{
		Code: constant.SUCCESS,
	}
	if len(meta) == 0 {
		return reply, errors.New("param length is zero")
	}
	message := buildPropertyMessage(meta, m)
	data, err := json.Marshal(message)
	if err != nil {
		return reply, nil
	}
	topic := buildPropertyTopic(m.EntityId, m.ModelId, m.PropertyType)
	if token := m.Client.Publish(topic, byte(0), false, data); token.WaitTimeout(5*time.Second) && token.Error() != nil {
		reply.Code = constant.FAIL
		reply.Data = token.Error().Error()
		return reply, nil
	}
	return reply, nil
}

// PubEvent 上报事件
func (m *MqttClient) PubEvent(ctx context.Context, meta define.PropertyKV, eventIdentifier string) (*define.Reply, error) {
	reply := &define.Reply{
		Code: constant.SUCCESS,
	}
	if len(meta) == 0 {
		return reply, errors.New("param length is zero")
	}

	message := buildEventMessage(meta, m, eventIdentifier)
	data, err := json.Marshal(message)
	if err != nil {
		return reply, nil
	}
	topic := buildEventTopic(m.EntityId, m.ModelId, eventIdentifier)
	// fmt.Printf("[PubEvent pub] topic:%s, message:%s\n", topic, string(data))
	if token := m.Client.Publish(topic, byte(0), false, data); token.WaitTimeout(5*time.Second) && token.Error() != nil {
		reply.Code = constant.FAIL
		reply.Data = token.Error().Error()
		return reply, nil
	}
	return reply, nil
}

// SubDeviceControl 同步订阅消息
func (m *MqttClient) SubDeviceControl(serviceIdentifier string) {
	topic := BuildServiceControlReply(m.ModelId, m.EntityId, serviceIdentifier)
	// fmt.Printf("[SubDeviceControl] topic:%s\n", topic)
	if token := m.Client.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
		// fmt.Printf("SubDeviceControl err:%s", token.Error())
	}

	<-m.UnSubScribeChan
	// fmt.Printf("[SubDeviceControl] closed, topic:%s\n", topic)
}

func (m *MqttClient) UnSubDeviceControl(serviceIdentifier string) error {
	// defer m.client.Disconnect(250)
	defer func() {
		close(m.UnSubScribeChan)
	}()
	topic := BuildServiceControlReply(m.ModelId, m.EntityId, serviceIdentifier)
	// fmt.Printf("[UnSubDeviceControl] topic:%s\n", topic)
	token := m.Client.Unsubscribe(topic)
	token.Wait()
	err := token.Error()
	if err != nil {
		return err
	}
	return nil
}
