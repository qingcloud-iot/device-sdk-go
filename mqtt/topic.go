package mqtt

import (
	"errors"
	"fmt"
	"git.internal.yunify.com/tools/device-sdk-go/index"
	"github.com/satori/go.uuid"
	"strings"
	"time"
)

/**
* @Author: hexing
* @Date: 19-9-9 下午2:37
 */
const (
	edge_topic = "/edge/%s/thing/event/+/post"
)
func buildAppPropertyTopic(appId string) string {
	return fmt.Sprintf("/edge/%s/thing/event/property/control", appId)
}
func buildAppEventTopic(appId, eventId string) string {
	return fmt.Sprintf("/edge/%s/thing/event/%s/control", appId, eventId)
}
func buildMessage(deviceId string, meta index.Metadata) *index.AppMessage {
	message := &index.AppMessage{
		Id:       uuid.NewV4().String(),
		Time:     time.Now().Unix(),
		DeviceId: deviceId,
		Params:   meta,
	}
	return message
}
func parseAppEvent(topic string) (string, string, error) {
	kv := strings.Split(topic, "/")
	if len(kv) != 7 {
		return "", "", errors.New("parse property error")
	}
	return kv[5], kv[2], nil
}
