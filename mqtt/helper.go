package mqtt

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"git.internal.yunify.com/iot-sdk/device-sdk-go/constant"
	"git.internal.yunify.com/iot-sdk/device-sdk-go/define"
	"github.com/dgrijalva/jwt-go"
	mqttp "github.com/eclipse/paho.mqtt.golang"

	uuid "github.com/satori/go.uuid"
)

const (
	QUIESCE = 30000 // milliseconds
)

// ParseToken return token payload
func ParseToken(deviceToken string) (string, string, error) {
	var (
		entityID string
		modelID  string
		err      error
	)
	defer func() {
		if errToken := recover(); err != nil {
			err = errToken.(error)
			return
		}
	}()
	token, _ := jwt.Parse(deviceToken, func(token *jwt.Token) (interface{}, error) {
		return "", nil
	})
	if payload, ok := token.Claims.(jwt.MapClaims); ok {
		if err := payload.Valid(); err != nil {
			return "", "", err
		}
		entityID, ok = payload[constant.TOKEN_DEVICE_ID].(string)
		if !ok {
			return "", "", errors.New("device id type error")
		}
		modelID, ok = payload[constant.TOKEN_THING_ID].(string)
		if !ok {
			return "", "", errors.New("device id type error")
		}
		return entityID, modelID, nil
	} else {
		return entityID, modelID, errors.New("token error")
	}
}

func buildPropertyMessage(meta define.PropertyKV, m *MqttClient) *define.ThingPropertyMsg {
	timeNow := time.Now().Unix() * 1000
	params := make(map[string]*define.PropertyValueAndTime)
	for k, v := range meta {
		property := &define.PropertyValueAndTime{
			Value: v,
			Time:  timeNow,
		}
		params[k] = property
	}
	message := &define.ThingPropertyMsg{
		ID:      uuid.NewV4().String(),
		Version: constant.MQTT_VERSION,
		Type:    constant.PROPERTY_TYPE,
		Params:  params,
		MetaData: define.MetaData{
			"modelId":   m.ModelId,
			"entityId":  m.EntityId,
			"epochTime": timeNow,
			"source":    []string{m.EntityId},
		},
	}
	return message
}

func buildPropertyMessageEx(meta define.PropertyKV, t int64) *define.ThingPropertyMsg {
	id := uuid.NewV4().String()
	params := make(map[string]*define.PropertyValueAndTime)
	for k, v := range meta {
		property := &define.PropertyValueAndTime{
			Value: v,
			Time:  t,
		}
		params[k] = property
	}
	message := &define.ThingPropertyMsg{
		ID:      id,
		Version: constant.MQTT_VERSION,
		Params:  params,
	}
	return message
}

func buildEventMessage(meta define.PropertyKV, m *MqttClient, eventIdentifier string) *define.ThingEventMsg {
	timeNow := time.Now().Unix() * 1000

	message := &define.ThingEventMsg{
		ID:      uuid.NewV4().String(),
		Version: constant.MQTT_VERSION,
		Type:    fmt.Sprintf("thing.event.%s.post", eventIdentifier),
		MetaData: define.MetaData{
			"modelId":   m.ModelId,
			"entityId":  m.EntityId,
			"epochTime": timeNow,
			"source":    []string{m.EntityId},
		},
		Params: &define.EventData{
			Value: meta,
			Time:  timeNow,
		},
	}
	return message
}

func ParseMessage(payload []byte) (*define.Message, error) {
	message := &define.Message{}
	err := json.Unmarshal(payload, message)
	if err != nil {
		// fmt.Errorf("parseMessage err:%s", err.Error())
		return nil, err
	}
	return message, nil
}
func buildPropertyTopic(entityID, modelID, propertyType string) string {
	return fmt.Sprintf(constant.POST_PROPERTY_TOPIC, modelID, entityID, propertyType)
}
func buildEventTopic(entityID, modelID, name string) string {
	return fmt.Sprintf(constant.POST_EVENT_TOPIC, modelID, entityID, name)
}

func buildPropertyReply(entityID, modelID string) string {
	return fmt.Sprintf(constant.SET_PROPERTY_TOPIC_REPLY, modelID, entityID)
}
func buildServiceReply(name, entityID, modelID string) string {
	return fmt.Sprintf(constant.SET_SERVICE_TOPIC_REPLY, modelID, entityID, name)
}

func BuildServiceControlReply(modelID, entityID, identifer string) string {
	return fmt.Sprintf(constant.DEVICE_CONTROL_TOPIC, modelID, entityID, identifer)
}

// Reply 服务调用的返回信息
func Reply(message *define.Message, client mqttp.Client, topic string, result define.PropertyKV) error {
	reply := &define.Reply{
		ID:   message.ID,
		Code: constant.SUCCESS,
		Data: make(define.PropertyKV),
	}
	reply.Data = result

	data, err := json.Marshal(reply)
	if err != nil {
		fmt.Printf("[recvDeviceControlReply] err:%s\n", err.Error())
		return err
	}
	fmt.Println(string(data))
	token := client.Publish(topic+"_reply", byte(0), false, data)
	if token.Error() != nil {
		fmt.Printf("[recvDeviceControlReply] err:%s\n", err.Error())
		return err
	}
	fmt.Printf("[recvDeviceControlReply] success\n")
	return nil
}
