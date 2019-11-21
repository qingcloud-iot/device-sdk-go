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
		opts.SetConnectionLostHandler(func(client mqttp.Client, e error) {
			fmt.Println("lost connect")
		})
		opts.SetOnConnectHandler(func(client mqttp.Client) {
			fmt.Println("connect success")
		})
		opts.SetDefaultPublishHandler(func(client mqttp.Client, msg mqttp.Message) {
			//fmt.Println("[sdk-go]", msg.Topic(), string(msg.Payload()))
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
		if token := client.Connect(); token.WaitTimeout(5*time.Second) && token.Error() != nil {
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
		m.cacheClient = cache.Cache(deviceId)
		m.pool = pool
		return m, nil
	}
}
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
	if c, ok := item.Data().(chan *index.Reply); ok {
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
		if err := m.client.Publish(topic+"_reply", byte(0), false, data); err != nil {
			fmt.Errorf("")
		}
	}
}
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
func (m *mqttClient) PubPropertyAsyncEx(meta index.Metadata, t int64) (index.ReplyChan, error) {
	ch := make(index.ReplyChan)
	reply := &index.Reply{
		Code: index.RPC_SUCCESS,
	}
	if len(meta) == 0 {
		return ch, errors.New("param length is zero")
	}
	message := buildPropertyMessageEx(meta, t)
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
func (m *mqttClient) PubEventSync(ctx context.Context, event string, meta index.Metadata) (*index.Reply, error) {
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
	fmt.Println("[PubEventSync] ", topic, string(data), time.Now().Unix())
	if token := m.client.Publish(topic, byte(0), false, data); token.WaitTimeout(5*time.Second) && token.Error() != nil {
		reply.Code = index.RPC_TIMEOUT
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
