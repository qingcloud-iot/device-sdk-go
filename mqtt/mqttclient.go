package mqtt

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"git.internal.yunify.com/iot-sdk/device-sdk-go/internal/client"
	"git.internal.yunify.com/iot-sdk/device-sdk-go/internal/constant"
	"git.internal.yunify.com/iot-sdk/device-sdk-go/internal/define"

	mqttp "github.com/eclipse/paho.mqtt.golang"
)

type DeviceControlHandler func(mqttp.Client, mqttp.Message)

type Options struct {
	Token         string               // 权限验证，及获取 ModelID、EntityID
	Server        string               // mqtt server
	PropertyType  string               // 属性分组（系统属性platform、基础属性base）
	MessageID     string               // 消息ID，设备内自增
	DeviceControl mqttp.MessageHandler // 设备控制的回调函数
	EntityID      string               // 设备 id
	ModelID       string               // 模型 id
}

type Client interface {
	mqttp.Client
}

type Message interface {
	Duplicate() bool
	Qos() byte
	Retained() bool
	Topic() string
	MessageID() uint16
	Payload() []byte
	Ack()
}

func (o *Options) SetDeviceControlHandler(deviceControl mqttp.MessageHandler) {
	o.DeviceControl = deviceControl
}

type MqttClient struct {
	Client mqttp.Client

	EntityId string
	ModelId  string

	PropertyType    string
	MessageID       string
	UnSubScribeChan chan bool
}

// NewMqtt 创建客户端实例
func NewMqtt(options *Options) (client.Client, error) {
	var (
		err error
	)

	m := &MqttClient{
		UnSubScribeChan: make(chan bool),
	}

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
	opts.SetDefaultPublishHandler(options.DeviceControl)

	client := mqttp.NewClient(opts)

	if err != nil {
		return nil, err
	}
	m.Client = client
	m.EntityId = options.EntityID
	m.ModelId = options.ModelID

	m.MessageID = options.MessageID
	if options.PropertyType != "" {
		m.PropertyType = options.PropertyType
	}
	return m, nil
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
		Code: constant.RPC_SUCCESS,
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
	// fmt.Printf("[PubProperty] topic:%s, message:%s\n", topic, string(data))
	if token := m.Client.Publish(topic, byte(0), false, data); token.WaitTimeout(5*time.Second) && token.Error() != nil {
		reply.Code = constant.RPC_TIMEOUT
		return reply, nil
	}
	return reply, nil
}

// PubEvent 上报事件
func (m *MqttClient) PubEvent(ctx context.Context, meta define.PropertyKV, eventIdentifier string) (*define.Reply, error) {
	reply := &define.Reply{
		Code: constant.RPC_SUCCESS,
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
		reply.Code = constant.RPC_TIMEOUT
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

// normal publish
// func (m *MqttClient) Publish(topic string, data []byte) (*index.Reply, error) {
// 	reply := &index.Reply{
// 		Code: index.RPC_SUCCESS,
// 	}
// 	if token := m.Client.Publish(topic+"_reply", byte(0), false, data); token.WaitTimeout(5*time.Second) && token.Error() != nil {
// 		fmt.Printf("[recvDeviceControlReply] err:%s", token.Error())
// 		reply.Code = index.RPC_FAIL
// 	} else {
// 		fmt.Println("[recvDeviceControlReply]", topic+"_reply", string(data))
// 	}

// 	return reply, nil
// }
