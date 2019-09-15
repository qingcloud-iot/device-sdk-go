package mqtt

import (
	"context"
	"encoding/json"
	"fmt"
	"git.internal.yunify.com/tools/device-sdk-go/index"
	mqttp "github.com/eclipse/paho.mqtt.golang"
	cache "github.com/muesli/cache2go"
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
		m := &mqttClient{
			client:      client,
			deviceId:    deviceId,
			thingId:     thingId,
			cacheClient: cache.Cache(deviceId),
		}
		return m, nil
	}
}
func (m *mqttClient) Start() error {
	if token := m.client.Subscribe(fmt.Sprintf(post_property_topic_reply, m.thingId, m.deviceId), byte(0), m.recvPropertyReply()); token.WaitTimeout(5*time.Second) && token.Error() != nil {
		return token.Error()
	}
	if token := m.client.Subscribe(fmt.Sprintf(post_event_topic_reply, m.thingId, m.deviceId), byte(0), m.recvEventReply()); token.WaitTimeout(5*time.Second) && token.Error() != nil {
		return token.Error()
	}
	if token := m.client.Subscribe(fmt.Sprintf(set_property_topic, m.thingId, m.deviceId), byte(0), m.setPropertyReply()); token.WaitTimeout(5*time.Second) && token.Error() != nil {
		return token.Error()
	}
	if token := m.client.Subscribe(fmt.Sprintf(set_service_topic, m.thingId, m.deviceId), byte(0), m.requestServiceReply()); token.WaitTimeout(5*time.Second) && token.Error() != nil {
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

	}
}
func (m *mqttClient) recvEventReply() func(mqttp.Client, mqttp.Message) {
	return func(client mqttp.Client, msg mqttp.Message) {
		topic := msg.Topic()
		//qos := msg.Qos()
		payload := msg.Payload()
		fmt.Println("[sdk-go-sub]", topic, string(payload))
	}
}
func (m *mqttClient) setPropertyReply() func(mqttp.Client, mqttp.Message) {
	return func(client mqtt.Client, msg mqtt.Message) {
		topic := msg.Topic()
		//qos := msg.Qos()
		payload := msg.Payload()
		fmt.Println("[sdk-go-sub]", topic, string(payload))
	}
}
func (m *mqttClient) requestServiceReply() func(mqttp.Client, mqttp.Message) {
	return func(client mqttp.Client, msg mqttp.Message) {
		topic := msg.Topic()
		//qos := msg.Qos()
		payload := msg.Payload()
		fmt.Println("[sdk-go-sub]", topic, string(payload))
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
		reply.Code = 1001
		reply.Message = ""
	}
	return reply
}
func pubPropertyReply(topic string, payload []byte) {

}
func (m *mqttClient) PubEvent(ctx context.Context, event string, meta index.Metadata) error {
	//message, err := json.Marshal(buildMessage(deviceId, meta))
	//if err != nil {
	//	return err
	//}
	//topic := buildAppEventTopic(m.name, event)
	//if token := m.client.Publish(topic, byte(0), false, message); token.WaitTimeout(5*time.Second) && token.Error() != nil {
	//	return token.Error()
	//}
	//fmt.Println("[sdk-go-pub]", topic, string(message))
	return nil
}
func (m *mqttClient) ReplyService(name string, meta index.Metadata) error {
	return nil
}
