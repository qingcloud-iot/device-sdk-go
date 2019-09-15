package index

/**
* @Author: hexing
* @Date: 19-9-9 上午11:32
 */
type MessageReply func(reply *Reply)
type ServiceHandle func(name string, meta Metadata)
type Options struct {
	DeviceId string
	Token    string
	Server   string
}
type Metadata map[string]*Property
type Property struct {
	Value interface{} `json:"value"`
	Time  int64       `json:"time"`
}
type Reply struct {
	Code    int
	Message string
	Data    interface{}
}
type Message struct {
	Id      string   `json:"id"`
	Version string   `json:"version"`
	Params  Metadata `json:"params"`
}
type Client interface {
	Start(messageReply MessageReply, serviceHandle ServiceHandle) error
	PubProperty(deviceId string, meta Metadata) error
	PubEvent(deviceId string, event string, meta Metadata) error
}
