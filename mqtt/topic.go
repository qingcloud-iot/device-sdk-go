package mqtt

import (
	"encoding/json"
	"fmt"
	"git.internal.yunify.com/tools/device-sdk-go/index"
	uuid "github.com/satori/go.uuid"
	"time"
)

/**
* @Author: hexing
* @Date: 19-9-9 下午2:37
 */
const (
	post_property_topic       = "/sys/%s/%s/thing/event/property/post"       //down
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
)

func buildPropertyMessage(meta index.Metadata) *index.ThingPropertyMsg {
	id := uuid.NewV4().String()
	params := make(map[string]*index.Property)
	for k, v := range meta {
		property := &index.Property{
			Value: v,
			Time:  time.Now().Unix() * 1000,
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
func buildPropertyMessageEx(meta index.Metadata, t int64) *index.ThingPropertyMsg {
	id := uuid.NewV4().String()
	params := make(map[string]*index.Property)
	for k, v := range meta {
		property := &index.Property{
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
func buildEventMessage(meta index.Metadata) *index.ThingEventMsg {
	id := uuid.NewV4().String()
	message := &index.ThingEventMsg{
		Id:      id,
		Version: MQTT_VERSION,
		Params: &index.EventData{
			Value: meta,
			Time:  time.Now().Unix() * 1000,
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
func buildProperty(deviceId, thingId string) string {
	return fmt.Sprintf(post_property_topic, thingId, deviceId)
}
func buildEvent(deviceId, thingId, name string) string {
	return fmt.Sprintf(post_event_topic, thingId, deviceId, name)
}

func buildPropertyReply(deviceId, thingId string) string {
	return fmt.Sprintf(set_property_topic_reply, thingId, deviceId)
}
func buildServiceReply(name, deviceId, thingId string) string {
	return fmt.Sprintf(set_service_topic_reply, thingId, deviceId, name)
}

func buildServiceControlReply(thingId, deviceId, identifer string) string {
	return fmt.Sprintf(device_control_topic, thingId, deviceId, identifer)
}
