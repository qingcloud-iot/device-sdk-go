package mqtt

/**
* @Author: hexing
* @Date: 19-9-9 下午2:37
 */
const (
	post_property_topic_reply = "/sys/%s/%s/thing/event/property/post_reply"//down
	post_event_topic_reply = "/sys/%s/%s//thing/event/+/post_reply" //down
	set_property_topic = "/sys/%s/%s/thing/service/property/set" //down
	set_service_topic = "/sys/%s/%s/thing/service/+"
	set_service_topic_reply = "/sys/%s/%s/thing/service/%s"
)

//func buildMessage(deviceId string, meta index.Metadata) *index.Message {
//	message := &index.Message{}
//	return message
//}
//func parseAppEvent(topic string) (string, string, error) {
//	kv := strings.Split(topic, "/")
//	if len(kv) != 7 {
//		return "", "", errors.New("parse property error")
//	}
//	return kv[5], kv[2], nil
//}
