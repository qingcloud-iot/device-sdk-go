package index

import "context"

type DownReply func(*Reply)
type SetProperty func(meta Metadata) (Metadata, error)
type ServiceHandle func(name string, meta Metadata) (Metadata, error)
type Options struct {
	Token     string // 权限验证，及获取thingid、deviceid
	Server    string // mqtt server
	Identifer string // sub 需定义
}
type ReplyChan chan *Reply
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
	Params  map[string]*Property `json:"params"`
}

//up event
type ThingEventMsg struct {
	Id      string     `json:"id"`
	Version string     `json:"version"`
	Params  *EventData `json:"params"`
}
type EventData struct {
	Value Metadata `json:"value"`
	Time  int64    `json:"time"`
}

type Reply struct {
	Code int         `json:"code"`
	Id   string      `json:"id"`
	Data interface{} `json:"data"`
}
type Message struct {
	Id      string   `json:"id"`
	Version string   `json:"version"`
	Params  Metadata `json:"params"`
}
type Client interface {
	Connect() error // 设备连接
	DisConnect()    // 设备断开连接
	//device
	PubProperty(ctx context.Context, meta Metadata) (*Reply, error) //post property sync
	PubPropertyAsync(meta Metadata) (ReplyChan, error)              //post property async
	//PubPropertyAsyncEx(meta Metadata, t int64) (ReplyChan, error)                  //post property async
	PubEvent(ctx context.Context, event string, meta Metadata) (*Reply, error) //post property　sync
	PubEventAsync(event string, meta Metadata) (ReplyChan, error)              //post property　async
	//driver
	PubTopicProperty(ctx context.Context, deviceId, thingId string, meta Metadata) (*Reply, error)            //post property sync
	PubTopicEvent(ctx context.Context, deviceId, thingId string, event string, meta Metadata) (*Reply, error) //post property　sync

	// sub
	SubDeviceControl()
	UnSubDeviceControl() error
}
