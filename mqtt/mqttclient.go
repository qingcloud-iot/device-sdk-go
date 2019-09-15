package mqtt

import (
	"encoding/json"
	"fmt"
	"git.internal.yunify.com/tools/device-sdk-go/index"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"strings"
	"time"
)

/**
* @Author: hexing
* @Date: 19-9-9 上午11:46
 */

type mqttClient struct {
	client mqtt.Client
	name   string
}

func NewMqtt(options *index.Options) (index.Client, error) {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(options.Server)
	opts.SetClientID("mqtt-" + options.Name)
	opts.SetUsername(options.Name)
	opts.SetPassword(options.Name)
	opts.SetCleanSession(true)
	opts.SetAutoReconnect(true)
	opts.SetKeepAlive(30 * time.Second)
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.WaitTimeout(5*time.Second) && token.Error() != nil {
		return nil, token.Error()
	}
	m := &mqttClient{
		client: client,
		name:   options.Name,
	}
	return m, nil
}
func (m *mqttClient) Start(propertyHandle index.PropertyHandle, eventHandle index.EventHandle) error {
	if token := m.client.Subscribe(fmt.Sprintf(edge_topic, m.name), byte(0), m.handle(propertyHandle, eventHandle)); token.WaitTimeout(5*time.Second) && token.Error() != nil {
		return token.Error()
	}
	return nil
}
func (m *mqttClient) PubProperty(deviceId string, meta index.Metadata) error {
	message, err := json.Marshal(buildMessage(deviceId, meta))
	if err != nil {
		return err
	}
	topic := buildAppPropertyTopic(m.name)
	if token := m.client.Publish(topic, byte(0), false, message); token.WaitTimeout(5*time.Second) && token.Error() != nil {
		return token.Error()
	}
	fmt.Println("[sdk-go-pub]", topic, string(message))
	return nil
}
func (m *mqttClient) PubEvent(deviceId string, event string, meta index.Metadata) error {
	message, err := json.Marshal(buildMessage(deviceId, meta))
	if err != nil {
		return err
	}
	topic := buildAppEventTopic(m.name, event)
	if token := m.client.Publish(topic, byte(0), false, message); token.WaitTimeout(5*time.Second) && token.Error() != nil {
		return token.Error()
	}
	fmt.Println("[sdk-go-pub]", topic, string(message))
	return nil
}
func (m *mqttClient) handle(propertyHandle index.PropertyHandle, eventHandle index.EventHandle) func(mqtt.Client, mqtt.Message) {
	return func(client mqtt.Client, msg mqtt.Message) {
		topic := msg.Topic()
		//qos := msg.Qos()
		payload := msg.Payload()
		fmt.Println("[sdk-go-sub]", topic, string(payload))
		switch {
		case strings.HasPrefix(topic, fmt.Sprintf("/edge/%s/thing/event/property/post", m.name)):
			message := &index.AppMessage{}
			err := json.Unmarshal(payload, message)
			if err != nil {
				fmt.Println(err)
				return
			}
			propertyHandle(message.DeviceId, message.Params)
		case strings.HasPrefix(topic, "/edge/") && strings.HasSuffix(topic, "/post"):
			message := &index.AppMessage{}
			err := json.Unmarshal(payload, message)
			if err != nil {
				fmt.Println(err)
				return
			}
			eventId, _, err := parseAppEvent(topic)
			if err != nil {
				fmt.Println(err)
				return
			}
			eventHandle(message.DeviceId, eventId, message.Params)
		default:
			fmt.Println(topic, string(payload))
			return
		}
	}
}
