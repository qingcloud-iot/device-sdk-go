package index

import "context"

const (
	PROPERTY_TYPE_BASE     = "base"
	PROPERTY_TYPE_PLATFORM = "platform"
)

const (
	PROPERTY_TYPE = "thing.property.post"
)

type DownReply func(*Reply)
type SetProperty func(meta PropertyKV) (PropertyKV, error)
type ServiceHandle func(name string, meta PropertyKV) (PropertyKV, error)
type Options struct {
	Token           string // 权限验证，及获取modelID、entityID
	Server          string // mqtt server
	PropertyType    string // 属性分组（系统属性platform、基础属性base）
	MessageID       string // 消息ID，设备内自增
	EventIdentifier string // 事件 identifier
	Identifer       string // sub 需定义
}
type ReplyChan chan *Reply
type PropertyKV map[string]interface{}
type PropertyValueAndTime struct {
	Value interface{} `json:"value"`
	Time  int64       `json:"time"`
}

type MetaData map[string]interface{}

type Event struct {
	Value PropertyKV `json:"value"`
	Time  int64      `json:"time"`
}

type ThingPropertyMsg struct {
	Id       string                           `json:"id"`
	Version  string                           `json:"version"`
	Params   map[string]*PropertyValueAndTime `json:"params"`
	Type     string                           `json:"type"`
	MetaData MetaData                         `json:"metaData"`
}

//up event
type ThingEventMsg struct {
	Id       string     `json:"id"`
	Version  string     `json:"version"`
	Params   *EventData `json:"params"`
	Type     string     `json:"type"`
	MetaData MetaData   `json:"metaData"`
}
type EventData struct {
	Value PropertyKV `json:"value"`
	Time  int64      `json:"time"`
}

type Reply struct {
	Code int         `json:"code"`
	Id   string      `json:"id"`
	Data interface{} `json:"data"`
}
type Message struct {
	Id      string     `json:"id"`
	Version string     `json:"version"`
	Type    string     `json:"type"`
	Params  PropertyKV `json:"params"`
}
type Client interface {
	Connect() error // 设备连接
	DisConnect()    // 设备断开连接
	//device
	PubProperty(ctx context.Context, meta PropertyKV) (*Reply, error) //post property sync
	PubPropertyAsync(meta PropertyKV) (ReplyChan, error)              //post property async
	//PubPropertyAsyncEx(meta Metadata, t int64) (ReplyChan, error)                  //post property async
	PubEvent(ctx context.Context, meta PropertyKV) (*Reply, error)  //post property　sync
	PubEventAsync(event string, meta PropertyKV) (ReplyChan, error) //post property　async
	//driver
	PubTopicProperty(ctx context.Context, entityID, modelID string, meta PropertyKV) (*Reply, error)            //post property sync
	PubTopicEvent(ctx context.Context, entityID, modelID string, event string, meta PropertyKV) (*Reply, error) //post property　sync

	// sub
	SubDeviceControl()
	UnSubDeviceControl() error
}
