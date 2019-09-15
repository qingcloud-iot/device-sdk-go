package index

import "context"

/**
* @Author: hexing
* @Date: 19-9-9 上午11:32
 */
type setProperty func(id string, meta Metadata)
type ServiceHandle func(id string, name string, meta Metadata)
type Options struct {
	DeviceId string
	Token    string
	Server   string
}
type Metadata map[string]interface{}
type Property struct {
	Value interface{} `json:"value"`
	Time  int64       `json:"time"`
}
type Reply struct {
	Code int
	Id   string
	Data interface{}
}
type Message struct {
	Id      string   `json:"id"`
	Version string   `json:"version"`
	Params  Metadata `json:"params"`
}
type Client interface {
	Start(setProperty setProperty, serviceHandle ServiceHandle) error
	PubProperty(ctx context.Context, meta Metadata) *Reply
	PubEvent(ctx context.Context, event string, meta Metadata) Reply
	ReplyProperty(id string, meta Metadata) error
	ReplyService(id string, name string, meta Metadata) error
}
