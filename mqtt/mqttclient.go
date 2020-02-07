package mqtt

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"git.internal.yunify.com/tools/device-sdk-go/index"
	mqttp "github.com/eclipse/paho.mqtt.golang"
	cache "github.com/muesli/cache2go"
	"github.com/panjf2000/ants"
	uuid "github.com/satori/go.uuid"
)

/**
* @Author: hexing
* @Date: 19-9-9 上午11:46
 */

type mqttClient struct {
	client      mqttp.Client
	deviceId    string
	thingId     string
	cacheClient *cache.CacheTable
	pool        *ants.Pool

	Identifier      string
	UnSubScribeChan chan bool
}

func NewHubMqtt(options *index.Options) (index.Client, error) {
	var (
		clientId string
		err      error
		server   string
	)
	if clientId == "" {
		clientId = "dirver-" + uuid.NewV4().String()
	}
	if server == "" {
		server = MQTT_HUB
	}
	m := &mqttClient{}
	opts := mqttp.NewClientOptions()
	opts.AddBroker(options.Server)
	opts.SetClientID(clientId)
	opts.SetUsername(clientId)
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

	client := mqttp.NewClient(opts)
	if token := client.Connect(); !token.WaitTimeout(5*time.Second) || token.Error() != nil {
		if token.Error() != nil {
			return nil, token.Error()
		}
		return m, fmt.Errorf("mqtt client connect fail")
	}
	if token := client.Subscribe(driver_set_service_topic, 0, func(client mqttp.Client, msg mqttp.Message) {
		m.requestServiceReply(options.ServiceHandle)(client, msg)
	}); token.Wait() && token.Error() != nil {
		return m, fmt.Errorf("mqtt client sub fail")
	}
	pool, err := ants.NewPool(WORKER_POOL)
	if err != nil {
		return nil, err
	}
	m.client = client
	m.pool = pool
	return m, nil
}

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
		return nil, err
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
		fmt.Printf("[sdk-go] topic:%s, message:%s\n", msg.Topic(), string(msg.Payload()))
		switch {
		case msg.Topic() == fmt.Sprintf(device_control_topic, thingId, deviceId, options.Identifer):
			m.recvDeviceControlReply(client, msg)
		default:
		}
	})
	client := mqttp.NewClient(opts)
	if token := client.Connect(); !token.WaitTimeout(5*time.Second) || token.Error() != nil {
		if token.Error() != nil {
			return nil, token.Error()
		}
		return m, fmt.Errorf("mqtt client connect fail")
	}
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

// 订阅到消息之后的回调
// 返回信息存储在 cache 里面
func (m *mqttClient) recvPropertyReply(client mqttp.Client, msg mqttp.Message) {
	topic := msg.Topic()
	//qos := msg.Qos()
	payload := msg.Payload()
	fmt.Println("[sdk-go-sub-property] reply", topic, string(payload), time.Now().Unix())
	reply := &index.Reply{}
	err := json.Unmarshal(payload, reply)
	if err != nil {
		fmt.Errorf("recvPropertyReply err:%s", err.Error())
		return
	}
	item, err := m.cacheClient.Value(reply.Id)
	if err != nil {
		fmt.Errorf("recvPropertyReply err:%s", err.Error())
		return
	}
	if c, ok := item.Data().(index.ReplyChan); ok {
		item.RemoveAboutToExpireCallback()
		if err := m.pool.Submit(func() {
			fmt.Println("[sdk-go-sub-property] reply success ", topic, string(payload))
			c <- reply
		}); err != nil {
			fmt.Errorf("pool exec err:%s", err.Error())
		}

	}
}

func (m *mqttClient) recvEventReply(client mqttp.Client, msg mqttp.Message) {
	topic := msg.Topic()
	//qos := msg.Qos()
	payload := msg.Payload()
	fmt.Println("[sdk-go-sub-event]", topic, string(payload))
	reply := &index.Reply{}
	err := json.Unmarshal(payload, reply)
	if err != nil {
		fmt.Errorf("recvPropertyReply err:%s", err.Error())
		return
	}
	item, err := m.cacheClient.Value(reply.Id)
	if err != nil {
		fmt.Errorf("recvPropertyReply err:%s", err.Error())
		return
	}
	ch := item.Data()
	if c, ok := ch.(chan *index.Reply); ok {
		item.RemoveAboutToExpireCallback()
		if err := m.pool.Submit(func() {
			c <- reply
		}); err != nil {
			fmt.Errorf("pool exec err:%s", err.Error())
		}

	}
}

//recv
func (m *mqttClient) setPropertyReply(setProperty index.SetProperty) func(mqttp.Client, mqttp.Message) {
	return func(client mqttp.Client, msg mqttp.Message) {
		var (
			topic   = msg.Topic()
			payload = msg.Payload()
			result  = index.Metadata{}
		)
		//qos := msg.Qos()
		fmt.Println("[sdk-go-set]", topic, string(payload))
		message, err := parseMessage(payload)
		if err != nil {
			fmt.Errorf("requestServiceReply err:%s", err.Error())
			return
		}
		reply := &index.Reply{
			Id:   message.Id,
			Code: index.RPC_SUCCESS,
			Data: make(index.Metadata),
		}
		if setProperty != nil {
			var err error
			result, err = setProperty(message.Params)
			if err != nil {
				reply.Code = index.RPC_FAIL
			}
		}
		reply.Data = result
		data, err := json.Marshal(reply)
		if err != nil {
			fmt.Errorf("requestServiceReply err:%s", err.Error())
			return
		}
		if token := m.client.Publish(topic+"_reply", byte(0), false, data); token.WaitTimeout(5*time.Second) && token.Error() != nil {
			fmt.Errorf("requestServiceReply err:%s", token.Error())
		} else {
			fmt.Println("[requestServiceReply]", topic+"_reply", string(data))
		}
	}
}

func (m *mqttClient) requestServiceReply(serviceHandle index.ServiceHandle) func(mqttp.Client, mqttp.Message) {
	return func(client mqttp.Client, msg mqttp.Message) {
		var (
			topic   = msg.Topic()
			payload = msg.Payload()
			result  = index.Metadata{}
		)
		name := parseServiceName(topic)
		//qos := msg.Qos()
		fmt.Println("[sdk-go-sub]", topic, string(payload))
		message, err := parseMessage(payload)
		if err != nil {
			fmt.Errorf("requestServiceReply err:%s", err.Error())
			return
		}
		reply := &index.Reply{
			Id:   message.Id,
			Code: index.RPC_SUCCESS,
			Data: make(index.Metadata),
		}
		if serviceHandle != nil {
			result, err = serviceHandle(name, message.Params)
			if err != nil {
				reply.Code = index.RPC_FAIL
			}
		}
		reply.Data = result
		data, err := json.Marshal(reply)
		if err != nil {
			fmt.Errorf("requestServiceReply err:%s", err.Error())
			return
		}
		if token := m.client.Publish(topic+"_reply", byte(0), false, data); token.WaitTimeout(5*time.Second) && token.Error() != nil {
			fmt.Errorf("requestServiceReply err:%s", token.Error())
		} else {
			fmt.Println("[requestServiceReply]", topic+"_reply", string(data))
		}
	}
}

// PubPropertySync 将消息 id 放入 cache 并设置过期时间，值为 chan reply，ctx 到期后返回
func (m *mqttClient) PubPropertySync(ctx context.Context, meta index.Metadata) (*index.Reply, error) {
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
	fmt.Println("[PubPropertySync] ", topic, string(data), time.Now().Unix())
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
	fmt.Println("[PubPropertyAsync] ", topic, string(data), time.Now().Unix())
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
func (m *mqttClient) PubEventSync(ctx context.Context, event string, meta index.Metadata) (*index.Reply, error) {
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
	fmt.Println("[PubEventSync] ", topic, string(data), time.Now().Unix())
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
func (m *mqttClient) PubTopicPropertySync(ctx context.Context, deviceId, thingId string, meta index.Metadata) (*index.Reply, error) {
	return nil, nil
}
func (m *mqttClient) PubTopicEventSync(ctx context.Context, deviceId, thingId string, event string, meta index.Metadata) (*index.Reply, error) {
	return nil, nil
}

// SubDeviceControlSync 同步订阅消息
func (m *mqttClient) SubDeviceControlSync() {
	topic := buildServiceControlReply(m.thingId, m.deviceId, m.Identifier)
	fmt.Println("SubDeviceControlSync topic:", topic)
	if token := m.client.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
		fmt.Printf("SubDeviceControlSync err:%s", token.Error())
	}

	<-m.UnSubScribeChan
	fmt.Printf("SubDeviceControlSync closed, topic:%s\n", topic)
}

func (m *mqttClient) UnSubDeviceControlSync() error {
	// defer m.client.Disconnect(250)
	defer func() {
		close(m.UnSubScribeChan)
	}()
	topic := buildServiceControlReply(m.thingId, m.deviceId, m.Identifier)
	fmt.Println("UnSubDeviceControlSync topic:", topic)
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
		fmt.Errorf("recvDeviceControlReply err:%s", err.Error())
		return
	}

	fmt.Printf("recvDeviceControlReply topic:%s payload:%s\n", topic, string(payload))

	reply := &index.Reply{
		Id:   message.Id,
		Code: index.RPC_SUCCESS,
		Data: make(index.Metadata),
	}

	reply.Data = message.Params

	data, err := json.Marshal(reply)
	if err != nil {
		fmt.Printf("recvDeviceControlReply err:%s\n", err.Error())
		return
	}
	if token := m.client.Publish(topic+"_reply", byte(0), false, data); token.WaitTimeout(5*time.Second) && token.Error() != nil {
		fmt.Printf("recvDeviceControlReply err:%s", token.Error())
	} else {
		fmt.Println("[recvDeviceControlReply]", topic+"_reply", string(data))
	}
}
