package mqtt

import (
	"encoding/json"
	"fmt"
	"git.internal.yunify.com/tools/device-sdk-go/index"
	"strconv"
	"time"
)

/**
* @Author: hexing
* @Date: 19-9-9 下午2:37
 */
const (
	post_property_topic       = "/sys/%s/%s/thing/event/property/post"        //down
	post_property_topic_reply = "/sys/%s/%s/thing/event/property/post_reply"  //down
	post_event_topic          = "/sys/%s/%s/thing/event/%s/post"              //down
	post_event_topic_reply    = "/sys/%s/%s//thing/event/+/post_reply"        //down
	set_property_topic        = "/sys/%s/%s/thing/service/property/set"       //down
	set_property_topic_reply  = "/sys/%s/%s/thing/service/property/set_reply" //down
	set_service_topic         = "/sys/%s/%s/thing/service/+"
	set_service_topic_reply   = "/sys/%s/%s/thing/service/%s"
)
const (
	RPC_TIME_OUT = 5 * time.Second
)

func buildMessage(meta index.Metadata) *index.Message {
	worker := GetInsIdWorker(WORKER_ID)
	id, _ := worker.NextId()
	str := strconv.FormatInt(id, 10)
	message := &index.Message{
		Id:      str,
		Version: MQTT_VERSION,
		Params:  meta,
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
func buildEvent(name, deviceId, thingId string) string {
	return fmt.Sprintf(post_event_topic, thingId, deviceId, name)
}
func buildPropertyReply(deviceId, thingId string) string {
	return fmt.Sprintf(set_property_topic_reply, thingId, deviceId)
}
func buildServiceReply(name, deviceId, thingId string) string {
	return fmt.Sprintf(set_service_topic_reply, thingId, deviceId, name)
}
