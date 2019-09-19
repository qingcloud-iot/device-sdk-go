package index

import "context"

/**
* @Author: hexing
* @Date: 19-9-9 上午11:32
 */
type SetProperty func(id string, meta Metadata)
type ServiceHandle func(id string, name string, meta Metadata)
type Options struct {
	DeviceId      string
	Token         string
	Server        string
	SetProperty   SetProperty
	ServiceHandle ServiceHandle
}
type Metadata map[string]interface{}
type Property struct {
	Value interface{} `json:"value"`
	Time  int64       `json:"time"`
}
type Event struct {
	Value Metadata `json:"value"`
	Time  int64    `json:"time"`
}
type ThingPropertyMsg struct {
	Id      string               `json:"id"`
	Version string               `json:"version"`
	Iotd    string               `json:"iotd"`
	Params  map[string]*Property `json:"params"`
}

//up event
type ThingEventMsg struct {
	Id      string     `json:"id"`
	Version string     `json:"version"`
	Iotd    string     `json:"iotd"`
	Params  *EventData `json:"params"`
}
type EventData struct {
	Value Metadata `json:"value"`
	Time  int64    `json:"time"`
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
	PubProperty(ctx context.Context, meta Metadata) (*Reply, error)            //post property sync
	PubEvent(ctx context.Context, event string, meta Metadata) (*Reply, error) //post property　sync
	ReplyProperty(reply *Reply) error
	ReplyService(name string, reply *Reply) error
}
