package mqtt

import (
	"api/metadata/v1"
	"context"
	"encoding/json"
	"fmt"
	"git.internal.yunify.com/tools/device-sdk-go/index"
	mqttp "github.com/eclipse/paho.mqtt.golang"
	cache "github.com/muesli/cache2go"
	"github.com/panjf2000/ants"
	"time"
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
}

func NewMqtt(options *index.Options) (index.Client, error) {
	if deviceId, thingId, err := parseToken(options.Token); err != nil {
		return nil, err
	} else {
		opts := mqttp.NewClientOptions()
		opts.AddBroker(options.Server)
		opts.SetClientID(deviceId)
		opts.SetUsername(deviceId)
		opts.SetPassword(options.Token)
		opts.SetCleanSession(true)
		opts.SetAutoReconnect(true)
		opts.SetKeepAlive(30 * time.Second)
		client := mqttp.NewClient(opts)
		if token := client.Connect(); token.WaitTimeout(5*time.Second) && token.Error() != nil {
			return nil, token.Error()
		}
		pool, err := ants.NewPool(WORKER_POOL)
		if err != nil {
			return nil, err
		}
		m := &mqttClient{
			client:      client,
			deviceId:    deviceId,
			thingId:     thingId,
			cacheClient: cache.Cache(deviceId),
			pool:        pool,
		}
		return m, nil
	}
}
func (m *mqttClient) Start(messageReply index.MessageReply, serviceHandle index.ServiceHandle) error {
	if token := m.client.Subscribe(fmt.Sprintf(post_property_topic_reply, m.thingId, m.deviceId), byte(0), m.recvPropertyReply()); token.WaitTimeout(5*time.Second) && token.Error() != nil {
		return token.Error()
	}
	if token := m.client.Subscribe(fmt.Sprintf(post_event_topic_reply, m.thingId, m.deviceId), byte(0), m.recvEventReply()); token.WaitTimeout(5*time.Second) && token.Error() != nil {
		return token.Error()
	}
	if token := m.client.Subscribe(fmt.Sprintf(set_property_topic, m.thingId, m.deviceId), byte(0), m.setPropertyReply()); token.WaitTimeout(5*time.Second) && token.Error() != nil {
		return token.Error()
	}
	if token := m.client.Subscribe(fmt.Sprintf(set_service_topic, m.thingId, m.deviceId), byte(0), m.requestServiceReply(serviceHandle)); token.WaitTimeout(5*time.Second) && token.Error() != nil {
		return token.Error()
	}
	return nil
}
func (m *mqttClient) recvPropertyReply() func(mqttp.Client, mqttp.Message) {
	return func(client mqttp.Client, msg mqttp.Message) {
		topic := msg.Topic()
		//qos := msg.Qos()
		payload := msg.Payload()
		fmt.Println("[sdk-go-sub]", topic, string(payload))
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
			if err := m.pool.Submit(func() {
				c <- reply
			}); err != nil {
				fmt.Errorf("pool exec err:%s", err.Error())
			}

		}

	}
}
func (m *mqttClient) recvEventReply() func(mqttp.Client, mqttp.Message) {
	return func(client mqttp.Client, msg mqttp.Message) {
		topic := msg.Topic()
		//qos := msg.Qos()
		payload := msg.Payload()
		fmt.Println("[sdk-go-sub]", topic, string(payload))
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
			if err := m.pool.Submit(func() {
				c <- reply
			}); err != nil {
				fmt.Errorf("pool exec err:%s", err.Error())
			}

		}
	}
}

//recv
func (m *mqttClient) setPropertyReply() func(mqttp.Client, mqttp.Message) {
	return func(client mqttp.Client, msg mqttp.Message) {
		topic := msg.Topic()
		//qos := msg.Qos()
		payload := msg.Payload()
		fmt.Println("[sdk-go-sub]", topic, string(payload))
	}
}
func (m *mqttClient) requestServiceReply(serviceHandle index.ServiceHandle) func(mqttp.Client, mqttp.Message) {
	return func(client mqttp.Client, msg mqttp.Message) {
		topic := msg.Topic()
		if isServceTopic(topic) {
			fmt.Errorf("requestServiceReply topic:%s is not match", topic)
			return
		}
		//name := parseServiceName(topic)
		//qos := msg.Qos()
		payload := msg.Payload()
		fmt.Println("[sdk-go-sub]", topic, string(payload))
		//message,err := parseMessage(payload)
		//if err != nil {
		//	fmt.Errorf("requestServiceReply err:%s",err.Error())
		//	return
		//}
		if serviceHandle != nil {
			//serviceHandle(message.Id,name,)
		}
	}
}
func (m *mqttClient) PubProperty(ctx context.Context, meta index.Metadata) *index.Reply {
	reply := &index.Reply{}
	message := buildMessage(meta)
	data, err := json.Marshal(message)
	if err != nil {
		return reply
	}
	topic := buildProperty(m.deviceId, m.thingId)
	if token := m.client.Publish(topic, byte(0), false, data); token.WaitTimeout(5*time.Second) && token.Error() != nil {
		return reply
	}
	ch := make(chan *index.Reply)
	m.cacheClient.Add(message.Id, RPC_TIME_OUT, reply)
	select {
	case value := <-ch:
		return value
	case <-ctx.Done():
		reply.Code = RPC_TIMEOUT
	}
	return reply
}
func (m *mqttClient) PubEvent(ctx context.Context, event string, meta index.Metadata) *index.Reply {
	reply := &index.Reply{}
	message := buildMessage(meta)
	data, err := json.Marshal(message)
	if err != nil {
		return reply
	}
	topic := buildEvent(m.deviceId, m.thingId, event)
	if token := m.client.Publish(topic, byte(0), false, data); token.WaitTimeout(5*time.Second) && token.Error() != nil {
		return reply
	}
	ch := make(chan *index.Reply)
	m.cacheClient.Add(message.Id, RPC_TIME_OUT, reply)
	select {
	case value := <-ch:
		return value
	case <-ctx.Done():
		reply.Code = RPC_TIMEOUT
	}
	return reply
}
func (m *mqttClient) ReplyProperty(id string, meta index.Metadata) error {
	return nil
}
func (m *mqttClient) ReplyService(id string, name string, meta index.Metadata) error {
	return nil
}
