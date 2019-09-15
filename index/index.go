package index

/**
* @Author: hexing
* @Date: 19-9-9 上午11:32
 */
type PropertyHandle func(meta Metadata)
type EventHandle func(event string, meta Metadata)
type Options struct {
	DeviceId string
	Token    string
	Server   string
}
type Metadata map[string]interface{}
type Client interface {
	Start(propertyHandle PropertyHandle, eventHandle EventHandle) error
	PubProperty(deviceId string, meta Metadata) error
	PubEvent(deviceId string, event string, meta Metadata) error
}
