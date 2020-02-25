package mqtt

import (
	"encoding/json"
	"fmt"
	"time"

	"git.internal.yunify.com/iot-sdk/device-sdk-go/index"
	uuid "github.com/satori/go.uuid"
)

const (
	post_property_topic       = "/sys/%s/%s/thing/property/%s/post"          // 属性上报 topic
	post_property_topic_reply = "/sys/%s/%s/thing/event/property/post_reply" //down
	post_event_topic          = "/sys/%s/%s/thing/event/%s/post"             //down
	post_topic_reply          = "/sys/%s/%s/thing/event/+/post_reply"        //down
	set_property_topic        = "/sys/%s/%s/thing/service/set/call"          //down
	set_property_topic_reply  = "/sys/%s/%s/thing/service/set/call_reply"    //down
	set_service_topic         = "/sys/%s/%s/thing/service/+/call"
	set_service_topic_reply   = "/sys/%s/%s/thing/service/%s/call_reply"

	device_control_topic = "/sys/%s/%s/thing/service/%s/call"
)
const (
	driver_set_service_topic = "/sys/+/+/thing/service/+/call"
)
const (
	MQTT_HUB     = "tcp://127.0.0.1:1883"
	RPC_TIME_OUT = 5 * time.Second
	QUIESCE      = 30000 // milliseconds
)

func buildPropertyMessage(meta index.PropertyKV, m *mqttClient) *index.ThingPropertyMsg {
	timeNow := time.Now().Unix() * 1000
	params := make(map[string]*index.PropertyValueAndTime)
	for k, v := range meta {
		property := &index.PropertyValueAndTime{
			Value: v,
			Time:  timeNow,
		}
		params[k] = property
	}
	message := &index.ThingPropertyMsg{
		Id:      m.MessageID,
		Version: MQTT_VERSION,
		Type:    index.PROPERTY_TYPE,
		Params:  params,
		MetaData: index.MetaData{
			"modelId":   m.ModelId,
			"entityId":  m.EntityId,
			"epochTime": timeNow,
			"source":    []string{m.EntityId},
		},
	}
	return message
}

func buildPropertyMessageEx(meta index.PropertyKV, t int64) *index.ThingPropertyMsg {
	id := uuid.NewV4().String()
	params := make(map[string]*index.PropertyValueAndTime)
	for k, v := range meta {
		property := &index.PropertyValueAndTime{
			Value: v,
			Time:  t,
		}
		params[k] = property
	}
	message := &index.ThingPropertyMsg{
		Id:      id,
		Version: MQTT_VERSION,
		Params:  params,
	}
	return message
}
func buildEventMessage(meta index.PropertyKV, m *mqttClient) *index.ThingEventMsg {
	timeNow := time.Now().Unix() * 1000

	message := &index.ThingEventMsg{
		Id:      m.MessageID,
		Version: MQTT_VERSION,
		Type:    fmt.Sprintf("thing.event.%s.post", m.EventIdentifier),
		MetaData: index.MetaData{
			"modelId":   m.ModelId,
			"entityId":  m.EntityId,
			"epochTime": timeNow,
			"source":    []string{m.EntityId},
		},
		Params: &index.EventData{
			Value: meta,
			Time:  timeNow,
		},
	}
	return message
}
func parseMessage(payload []byte) (*index.Message, error) {
	message := &index.Message{}
	err := json.Unmarshal(payload, message)
	if err != nil {
		fmt.Errorf("parseMessage err:%s", err.Error())
		return nil, err
	}
	return message, nil
}
func buildPropertyTopic(entityID, modelID, propertyType string) string {
	return fmt.Sprintf(post_property_topic, modelID, entityID, propertyType)
}
func buildEventTopic(entityID, modelID, name string) string {
	return fmt.Sprintf(post_event_topic, modelID, entityID, name)
}

func buildPropertyReply(entityID, modelID string) string {
	return fmt.Sprintf(set_property_topic_reply, modelID, entityID)
}
func buildServiceReply(name, entityID, modelID string) string {
	return fmt.Sprintf(set_service_topic_reply, modelID, entityID, name)
}

func buildServiceControlReply(modelID, entityID, identifer string) string {
	return fmt.Sprintf(device_control_topic, modelID, entityID, identifer)
}
