package mqtt

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/qingcloud-iot/device-sdk-go/register"

	iClient "github.com/qingcloud-iot/device-sdk-go/client"
	"github.com/qingcloud-iot/device-sdk-go/constant"
	"github.com/qingcloud-iot/device-sdk-go/define"

	mqttp "github.com/eclipse/paho.mqtt.golang"
)

const (
	DefaultWaitTimeout = 3 * time.Second
)

//
type Handler func(msg *define.Message) define.PropertyKV

// DeviceControlHandler 服务调用结构体，用于处理下行数据的业务逻辑
type DeviceControlHandler struct {
	ServiceIdentifer string
	ServiceHandler   iClient.CallBack
}

type Options struct {
	Token string // 权限验证，及获取 ModelID、EntityID

	TLS bool // 使用 tls 连接

	MiddleCredential       string // 批量设备注册的中间凭证
	DynamocRegisterAddress string // 动态注册的服务地址

	Server         string                 // mqtt server
	PropertyType   string                 // 属性分组（系统属性platform、基础属性base）
	DeviceHandlers []DeviceControlHandler // 服务调用的回调函数
	EntityID       string                 // 设备 id
	ModelID        string                 // 模型 id

	CertFilePath string // mqtts 证书地址

	AutoReconnect   bool      // 是否自动重连
	LostConnectChan chan bool // 掉线后通知 chan，方便程序自处理
	ReConnectChan   chan bool // 掉线后重连 chan，通知设备已经重连
	n               int

	KeepAlive            time.Duration // 心跳间隔, 默认 30s
	MaxReconnectInterval time.Duration // 最大重连间隔, 默认 10 * time.Minute
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

	mqttp.CRITICAL = log.New(os.Stderr, "", log.Lshortfile|log.LstdFlags)
	mqttp.ERROR = log.New(os.Stderr, "", log.Lshortfile|log.LstdFlags)
	mqttp.WARN = log.New(os.Stderr, "", log.Lshortfile|log.LstdFlags)
	mqttp.DEBUG = log.New(os.Stdout, "", log.Lshortfile|log.LstdFlags)

	// use mqtts communicate
	if options.TLS {
		log.Println("use mqtts!")
		tlsConfig := &tls.Config{}
		if options.CertFilePath != "" {
			cert := x509.NewCertPool()
			pemCerts, err := ioutil.ReadFile(options.CertFilePath)
			if err == nil {
				if !cert.AppendCertsFromPEM(pemCerts) {
					log.Panic("failed to parse root certificate")
				}
			} else {
				log.Panic(err)
			}
			tlsConfig.RootCAs = cert
		}

		opts.SetTLSConfig(tlsConfig)
		opts.AddBroker("ssl://" + options.Server)
	} else {
		opts.AddBroker("tcp://" + options.Server)
	}
	if options.EntityID != "" {
		opts.SetClientID(options.EntityID)
		opts.SetUsername(options.EntityID)
	}
	opts.SetPassword(options.Token)
	opts.SetCleanSession(true)
	if options.AutoReconnect {
		log.Println("This device will auto reconnect to ehub/ihub!")
		opts.SetAutoReconnect(true)
	} else {
		log.Println("This device will not auto reconnect to ehub/ihub, you can ensure reconnect by set the config param <auto_reconnect:true>!")
		opts.SetAutoReconnect(false)
	}
	if options.KeepAlive != 0 {
		opts.SetKeepAlive(options.KeepAlive)
	}
	if options.MaxReconnectInterval != 0 {
		opts.SetMaxReconnectInterval(options.MaxReconnectInterval)
	}
	opts.SetConnectionLostHandler(func(client mqttp.Client, e error) {
		if options.LostConnectChan != nil {
			options.LostConnectChan <- true
		}
	})
	opts.SetOnConnectHandler(func(client mqttp.Client) {
		options.n++
		// 表示重连
		if options.n != 1 && options.ReConnectChan != nil {
			log.Println("reconnect ehub/ihub success!")
			options.ReConnectChan <- true
		} else {
			log.Println("connect ehub/ihub success!")
		}
	})

	if options.DeviceHandlers != nil {
		opts.SetDefaultPublishHandler(func(client mqttp.Client, msg mqttp.Message) {
			log.Printf("[sdk-go sub] topic: %s, paload: %s\n", msg.Topic(), string(msg.Payload()))

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
					result := handler.ServiceHandler.Handler(message)

					// reply
					if err = Reply(message, client, topic, result); err != nil {
						log.Printf("topic:%s, reply error: %s\n", topic, err.Error())
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
	if token := m.Client.Connect(); !token.WaitTimeout(5 * time.Second) || token.Error() != nil {
		if token.Error() != nil {
			return fmt.Errorf("连接错误:%s", token.Error())
		}
		return errors.New("连接超时，请检查网络及配置")
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

// PubPropertyWithTime 上报自定义时间属性
func (m *MqttClient) PubPropertyWithTime(ctx context.Context, metaDataWithTime define.PropertyKVWithTime) (*define.Reply, error) {

	reply := &define.Reply{
		Code: constant.SUCCESS,
	}
	if len(metaDataWithTime) == 0 {
		return reply, errors.New("param length is zero")
	}
	message := buildPropertyMessageWithTIme(metaDataWithTime, m)
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
	// log.Printf("[PubEvent pub] topic:%s, message:%s\n", topic, string(data))
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
	// log.Printf("[SubDeviceControl] topic:%s\n", topic)
	if token := m.Client.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
		// log.Printf("SubDeviceControl err:%s", token.Error())
	}

	<-m.UnSubScribeChan
	// log.Printf("[SubDeviceControl] closed, topic:%s\n", topic)
}

func (m *MqttClient) UnSubDeviceControl(serviceIdentifier string) error {
	// defer m.client.Disconnect(250)
	defer func() {
		close(m.UnSubScribeChan)
	}()
	topic := BuildServiceControlReply(m.ModelId, m.EntityId, serviceIdentifier)
	// log.Printf("[UnSubDeviceControl] topic:%s\n", topic)
	token := m.Client.Unsubscribe(topic)
	token.Wait()
	err := token.Error()
	if err != nil {
		return err
	}
	return nil
}

func (m *MqttClient) Subscribe(topic string, qos int32, cb iClient.MessageCallback) error {
	if topic == "" || qos < 0 || qos > 2 || cb == nil {
		return errors.New("invalid arguments")
	}
	if token := m.Client.Subscribe(topic, byte(qos), func(client mqttp.Client, message mqttp.Message) {
		cb(message.Topic(), message.Payload())
	}); token.WaitTimeout(DefaultWaitTimeout) && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (m *MqttClient) SubscribeMultiple(topics []string, cb iClient.MessageCallback) error {
	filters := make(map[string]byte)
	for _, topic := range topics {
		if topic == "" {
			continue
		}
		filters[topic] = 0
	}
	if token := m.Client.SubscribeMultiple(filters, func(client mqttp.Client, message mqttp.Message) {
		cb(message.Topic(), message.Payload())
	}); token.WaitTimeout(DefaultWaitTimeout) && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (m *MqttClient) Unsubscribe(topics []string) error {
	if token := m.Client.Unsubscribe(topics...); token.WaitTimeout(DefaultWaitTimeout) && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (m *MqttClient) Publish(topic string, qos int32, payload []byte) error {
	if topic == "" || qos < 0 || qos > 2 || payload == nil {
		return errors.New("invalid arguments")
	}
	if token := m.Client.Publish(topic, byte(qos), false, payload); token.WaitTimeout(DefaultWaitTimeout) && token.Error() != nil {
		return token.Error()
	}
	return nil
}
