package mqtt

import (
	"context"
	"encoding/json"
	"errors"
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
	m := &mqttClient{}
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
		opts.SetDefaultPublishHandler(func(client mqttp.Client, msg mqttp.Message) {
			fmt.Println("[sdk-go]", msg.Topic(), string(msg.Payload()))
			switch {
			case msg.Topic() == fmt.Sprintf(post_property_topic_reply, thingId, deviceId):
				m.recvPropertyReply(client, msg)
			case msg.Topic() == fmt.Sprintf(set_property_topic, thingId, deviceId):
				m.setPropertyReply(options.SetProperty)(client, msg)
			case isServiceTopic(thingId, deviceId, msg.Topic()):
				m.requestServiceReply(options.ServiceHandle)(client, msg)
			default:
				m.recvEventReply(client, msg)
			}
		})
		client := mqttp.NewClient(opts)
		if token := client.Connect(); token.Wait() && token.Error() != nil {
			return nil, token.Error()
		}
		pool, err := ants.NewPool(WORKER_POOL)
		if err != nil {
			return nil, err
		}
		m.client = client
		m.deviceId = deviceId
		m.thingId = thingId
		m.cacheClient = cache.Cache(deviceId)
		m.pool = pool
		return m, nil
	}
}
func (m *mqttClient) recvPropertyReply(client mqttp.Client, msg mqttp.Message) {
	topic := msg.Topic()
	//qos := msg.Qos()
	payload := msg.Payload()
	fmt.Println("[sdk-go-sub-property]", topic, string(payload))
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
		if err := m.client.Publish(topic+"_reply", byte(0), false, data); err != nil {
			fmt.Errorf("")
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
		if err := m.client.Publish(topic+"_reply", byte(0), false, data); err != nil {
			fmt.Errorf("")
		}
	}
}
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
	if token := m.client.Publish(topic, byte(0), false, data); token.WaitTimeout(5*time.Second) && token.Error() != nil {
		reply.Code = index.RPC_PUBLISH_TIMEOUT
		return reply, nil
	}
	fmt.Println(topic, string(data))
	ch := make(chan *index.Reply)
	m.cacheClient.Add(message.Id, RPC_TIME_OUT, ch)
	select {
	case value := <-ch:
		return value, nil
	case <-ctx.Done():
		reply.Code = index.RPC_TIMEOUT
	}
	return reply, nil
}
func (m *mqttClient) PubEvent(ctx context.Context, event string, meta index.Metadata) (*index.Reply, error) {
	reply := &index.Reply{
		Code: index.RPC_SUCCESS,
	}
	if len(meta) == 0 {
		return reply, errors.New("param length is zero")
	}
	message := buildEventMessage(meta)
	data, err := json.Marshal(message)
	if err != nil {
		return reply, nil
	}
	topic := buildEvent(m.deviceId, m.thingId, event)
	fmt.Println(topic, string(data))
	if token := m.client.Publish(topic, byte(0), false, data); token.WaitTimeout(5*time.Second) && token.Error() != nil {
		reply.Code = index.RPC_PUBLISH_TIMEOUT
		return reply, nil
	}
	ch := make(chan *index.Reply)
	m.cacheClient.Add(message.Id, RPC_TIME_OUT, ch)
	select {
	case value := <-ch:
		return value, nil
	case <-ctx.Done():
		reply.Code = index.RPC_TIMEOUT
	}
	return reply, nil
}

//func (m *mqttClient) HandleProperty(reply *index.Reply) error {
//	topic := buildPropertyReply(m.deviceId, m.thingId)
//	data, err := json.Marshal(reply)
//	if err != nil {
//		return err
//	}
//	if token := m.client.Publish(topic, byte(0), false, data); token.WaitTimeout(5*time.Second) && token.Error() != nil {
//		return token.Error()
//	}
//	return nil
//}
//func (m *mqttClient) HandleService(name string, reply *index.Reply) error {
//	topic := buildServiceReply(name, m.deviceId, m.thingId)
//	data, err := json.Marshal(reply)
//	if err != nil {
//		return err
//	}
//	if token := m.client.Publish(topic, byte(0), false, data); token.WaitTimeout(5*time.Second) && token.Error() != nil {
//		return token.Error()
//	}
//	return nil
//}
