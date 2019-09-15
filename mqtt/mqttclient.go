package mqtt

import (
	"context"
	"fmt"
	"git.internal.yunify.com/tools/device-sdk-go/index"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"time"
)

/**
* @Author: hexing
* @Date: 19-9-9 上午11:46
 */

type mqttClient struct {
	client mqtt.Client
	deviceId string
	thingId string
}

func NewMqtt(options *index.Options) (index.Client, error) {
	if deviceId,thingId,err := parseToken(options.Token);err != nil {
		return nil,err
	}else {
		opts := mqtt.NewClientOptions()
		opts.AddBroker(options.Server)
		opts.SetClientID(deviceId)
		opts.SetUsername(deviceId)
		opts.SetPassword(options.Token)
		opts.SetCleanSession(true)
		opts.SetAutoReconnect(true)
		opts.SetKeepAlive(30 * time.Second)
		client := mqtt.NewClient(opts)
		if token := client.Connect(); token.WaitTimeout(5*time.Second) && token.Error() != nil {
			return nil, token.Error()
		}
		m := &mqttClient{
			client: client,
			deviceId:deviceId,
			thingId:thingId,
		}
		return m, nil
	}
}
func (m *mqttClient) Start(propertyHandle index.PropertyHandle, eventHandle index.EventHandle) error {
	if token := m.client.Subscribe(fmt.Sprintf(post_property_topic_reply, m.thingId,m.deviceId), byte(0), m.recvPropertyReply()); token.WaitTimeout(5*time.Second) && token.Error() != nil {
		return token.Error()
	}
	if token := m.client.Subscribe(fmt.Sprintf(post_event_topic_reply, m.thingId,m.deviceId), byte(0), m.recvEventReply()); token.WaitTimeout(5*time.Second) && token.Error() != nil {
		return token.Error()
	}
	if token := m.client.Subscribe(fmt.Sprintf(set_property_topic, m.thingId,m.deviceId), byte(0), m.setPropertyReply()); token.WaitTimeout(5*time.Second) && token.Error() != nil {
		return token.Error()
	}
	if token := m.client.Subscribe(fmt.Sprintf(set_service_topic, m.thingId,m.deviceId), byte(0), m.requestServiceReply()); token.WaitTimeout(5*time.Second) && token.Error() != nil {
		return token.Error()
	}
	return nil
}
func (m *mqttClient)recvPropertyReply() func(mqtt.Client, mqtt.Message) {
	return func(client mqtt.Client, msg mqtt.Message) {
		topic := msg.Topic()
		//qos := msg.Qos()
		payload := msg.Payload()
		fmt.Println("[sdk-go-sub]", topic, string(payload))
	}
}
func (m *mqttClient)recvEventReply() func(mqtt.Client, mqtt.Message) {
	return func(client mqtt.Client, msg mqtt.Message) {
		topic := msg.Topic()
		//qos := msg.Qos()
		payload := msg.Payload()
		fmt.Println("[sdk-go-sub]", topic, string(payload))
	}
}
func (m *mqttClient)setPropertyReply() func(mqtt.Client, mqtt.Message) {
	return func(client mqtt.Client, msg mqtt.Message) {
		topic := msg.Topic()
		//qos := msg.Qos()
		payload := msg.Payload()
		fmt.Println("[sdk-go-sub]", topic, string(payload))
	}
}
func (m *mqttClient)requestServiceReply() func(mqtt.Client, mqtt.Message) {
	return func(client mqtt.Client, msg mqtt.Message) {
		topic := msg.Topic()
		//qos := msg.Qos()
		payload := msg.Payload()
		fmt.Println("[sdk-go-sub]", topic, string(payload))
	}
}
func (m *mqttClient) PubProperty(ctx context.Context,meta index.Metadata) error {
	//message, err := json.Marshal(buildMessage(meta))
	//if err != nil {
	//	return err
	//}
	//topic := buildAppPropertyTopic(m.name)
	//if token := m.client.Publish(topic, byte(0), false, message); token.WaitTimeout(5*time.Second) && token.Error() != nil {
	//	return token.Error()
	//}
	//fmt.Println("[sdk-go-pub]", topic, string(message))
	return nil
}
func pubPropertyReply(topic string,payload []byte) {

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
func (m *mqttClient) ReplyService(name string,meta index.Metadata) error {
	return nil
}