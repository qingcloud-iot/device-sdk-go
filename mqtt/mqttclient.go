package mqtt

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"git.internal.yunify.com/iot-sdk/device-sdk-go/index"
	mqttp "github.com/eclipse/paho.mqtt.golang"
	cache "github.com/muesli/cache2go"
	"github.com/panjf2000/ants"
)

type mqttClient struct {
	client      mqttp.Client
	deviceId    string
	thingId     string
	cacheClient *cache.CacheTable
	pool        *ants.Pool

	Identifier      string
	UnSubScribeChan chan bool
}

// NewMqtt 创建客户端实例
func NewMqtt(options *index.Options) (index.Client, error) {
	var (
		deviceId string
		thingId  string
		err      error
	)
	m := &mqttClient{
		UnSubScribeChan: make(chan bool),
	}
	if deviceId, thingId, err = parseToken(options.Token); err != nil {
		return nil, errors.New("Parse token error: " + err.Error())
	}
	opts := mqttp.NewClientOptions()
	opts.AddBroker(options.Server)
	opts.SetClientID(deviceId)
	opts.SetUsername(deviceId)
	opts.SetPassword(options.Token)
	opts.SetCleanSession(true)
	opts.SetAutoReconnect(true)
	opts.SetKeepAlive(30 * time.Second)
	opts.SetConnectionLostHandler(func(client mqttp.Client, e error) {
		fmt.Println("lost connect")
	})
	opts.SetOnConnectHandler(func(client mqttp.Client) {
		fmt.Println("connect success")
	})
	opts.SetDefaultPublishHandler(func(client mqttp.Client, msg mqttp.Message) {
		fmt.Printf("[sdk-go sub] topic:%s, message:%s\n", msg.Topic(), string(msg.Payload()))
		switch {
		case msg.Topic() == fmt.Sprintf(device_control_topic, thingId, deviceId, options.Identifer):
			m.recvDeviceControlReply(client, msg)
		default:
		}
	})
	client := mqttp.NewClient(opts)

	pool, err := ants.NewPool(WORKER_POOL)
	if err != nil {
		return nil, err
	}
	m.client = client
	m.deviceId = deviceId
	m.thingId = thingId
	m.Identifier = options.Identifer
	m.cacheClient = cache.Cache(deviceId)
	m.pool = pool
	return m, nil
}

// Connect 连接 ihub 或 ehub
func (m *mqttClient) Connect() error {
	if token := m.client.Connect(); !token.WaitTimeout(5*time.Second) || token.Error() != nil {
		return fmt.Errorf("mqtt client connect fail")
	}
	return nil
}

// PubProperty 将消息 id 放入 cache 并设置过期时间，值为 chan reply，ctx 到期后返回
func (m *mqttClient) PubProperty(ctx context.Context, meta index.Metadata) (*index.Reply, error) {
	reply := &index.Reply{
		Code: index.RPC_SUCCESS,
	}
	if len(meta) == 0 {
		return reply, errors.New("param length is zero")
	}
	message := buildPropertyMessage(meta)
	data, err := json.Marshal(message)
	if err != nil {
		return reply, nil
	}
	topic := buildProperty(m.deviceId, m.thingId)
	fmt.Printf("[PubProperty] topic:%s, message:%s\n", topic, string(data))
	if token := m.client.Publish(topic, byte(0), false, data); token.WaitTimeout(5*time.Second) && token.Error() != nil {
		reply.Code = index.RPC_TIMEOUT
		return reply, nil
	}
	return reply, nil
}

func (m *mqttClient) PubPropertyAsync(meta index.Metadata) (index.ReplyChan, error) {
	ch := make(index.ReplyChan)
	reply := &index.Reply{
		Code: index.RPC_SUCCESS,
	}
	if len(meta) == 0 {
		return ch, errors.New("param length is zero")
	}
	message := buildPropertyMessage(meta)
	data, err := json.Marshal(message)
	if err != nil {
		return ch, err
	}
	topic := buildProperty(m.deviceId, m.thingId)
	fmt.Printf("[PubPropertyAsync] topic:%s, message:%s\n", topic, string(data))
	if token := m.client.Publish(topic, byte(0), false, data); token.WaitTimeout(5*time.Second) && token.Error() != nil {
		reply.Code = index.RPC_TIMEOUT
		return ch, token.Error()
	}
	item := m.cacheClient.Add(message.Id, RPC_TIME_OUT, ch)
	item.SetAboutToExpireCallback(func(i interface{}) {
		fmt.Printf("[PubPropertyAsync] i:%+v,timeout topic:%s,data:%s", i, topic, string(data))
		reply := &index.Reply{
			Code: index.RPC_TIMEOUT,
		}
		ch <- reply
	})
	return ch, nil
}

// PubEventSync event 就是将整个 meta 放到 中
func (m *mqttClient) PubEvent(ctx context.Context, event string, meta index.Metadata) (*index.Reply, error) {
	reply := &index.Reply{
		Code: index.RPC_SUCCESS,
	}
	if len(meta) == 0 {
		return reply, errors.New("param length is zero")
	}

	/*
		message{
				id: uuid
				version: v1.0.0
				params: EventData{
					Value: meta,
					Time:  time.Now().Unix() * 1000,
				}
			}
	*/
	message := buildEventMessage(meta)
	data, err := json.Marshal(message)
	if err != nil {
		return reply, nil
	}
	topic := buildEvent(m.deviceId, m.thingId, event)
	fmt.Printf("[PubEvent pub] topic:%s, message:%s\n", topic, string(data))
	if token := m.client.Publish(topic, byte(0), false, data); token.WaitTimeout(5*time.Second) && token.Error() != nil {
		reply.Code = index.RPC_TIMEOUT
		return reply, nil
	}
	return reply, nil
}

func (m *mqttClient) PubEventAsync(event string, meta index.Metadata) (index.ReplyChan, error) {
	ch := make(index.ReplyChan)
	if len(meta) == 0 {
		return ch, errors.New("param length is zero")
	}
	message := buildEventMessage(meta)
	data, err := json.Marshal(message)
	if err != nil {
		return ch, err
	}
	topic := buildEvent(m.deviceId, m.thingId, event)
	fmt.Println(topic, string(data))
	if token := m.client.Publish(topic, byte(0), false, data); token.WaitTimeout(5*time.Second) && token.Error() != nil {
		return ch, err
	}
	item := m.cacheClient.Add(message.Id, RPC_TIME_OUT, ch)
	item.AddAboutToExpireCallback(func(i interface{}) {
		reply := &index.Reply{
			Code: index.RPC_TIMEOUT,
		}
		ch <- reply
	})
	return ch, nil
}

//driver
func (m *mqttClient) PubTopicProperty(ctx context.Context, deviceId, thingId string, meta index.Metadata) (*index.Reply, error) {
	return nil, nil
}
func (m *mqttClient) PubTopicEvent(ctx context.Context, deviceId, thingId string, event string, meta index.Metadata) (*index.Reply, error) {
	return nil, nil
}

// SubDeviceControl 同步订阅消息
func (m *mqttClient) SubDeviceControl() {
	topic := buildServiceControlReply(m.thingId, m.deviceId, m.Identifier)
	fmt.Printf("[SubDeviceControl] topic:%s\n", topic)
	if token := m.client.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
		fmt.Printf("SubDeviceControl err:%s", token.Error())
	}

	<-m.UnSubScribeChan
	fmt.Printf("[SubDeviceControl] closed, topic:%s\n", topic)
}

func (m *mqttClient) UnSubDeviceControl() error {
	// defer m.client.Disconnect(250)
	defer func() {
		close(m.UnSubScribeChan)
	}()
	topic := buildServiceControlReply(m.thingId, m.deviceId, m.Identifier)
	fmt.Printf("[UnSubDeviceControl] topic:%s\n", topic)
	token := m.client.Unsubscribe(topic)
	token.Wait()
	err := token.Error()
	if err != nil {
		return err
	}
	return nil
}

// recvDeviceControlReply 订阅消息后的回调函数
func (m *mqttClient) recvDeviceControlReply(client mqttp.Client, msg mqttp.Message) {

	topic := msg.Topic()
	payload := msg.Payload()

	//qos := msg.Qos()
	fmt.Println("[sdk-go-device-control] ", topic, string(payload))
	message, err := parseMessage(payload)
	if err != nil {
		fmt.Printf("recvDeviceControlReply err:%s", err.Error())
		return
	}

	fmt.Printf("[recvDeviceControlReply] topic:%s payload:%s\n", topic, string(payload))

	reply := &index.Reply{
		Id:   message.Id,
		Code: index.RPC_SUCCESS,
		Data: make(index.Metadata),
	}

	reply.Data = message.Params

	data, err := json.Marshal(reply)
	if err != nil {
		fmt.Printf("[recvDeviceControlReply] err:%s\n", err.Error())
		return
	}
	if token := m.client.Publish(topic+"_reply", byte(0), false, data); token.WaitTimeout(5*time.Second) && token.Error() != nil {
		fmt.Printf("[recvDeviceControlReply] err:%s", token.Error())
	} else {
		fmt.Println("[recvDeviceControlReply]", topic+"_reply", string(data))
	}
}
