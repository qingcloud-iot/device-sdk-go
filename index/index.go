package index

import (
	"context"
)

const (
	// 上报系统属性时，用于构造 topic
	PROPERTY_TYPE_BASE = "base"

	// 上报设备属性时，用于构造 topic
	PROPERTY_TYPE_PLATFORM = "platform"

	// 属性上报时的 type 字段值
	PROPERTY_TYPE = "thing.property.post"
)

// PropertyKV 事件 kv 健值对
type PropertyKV map[string]interface{}

// PropertyValueAndTime 属性
type PropertyValueAndTime struct {
	Value interface{} `json:"value"`
	Time  int64       `json:"time"`
}

// MetaData 上报元数据
type MetaData map[string]interface{}

// ThingPropertyMsg 上报属性的数据结构
type ThingPropertyMsg struct {
	ID       string                           `json:"id"`
	Version  string                           `json:"version"`
	Params   map[string]*PropertyValueAndTime `json:"params"`
	Type     string                           `json:"type"`
	MetaData MetaData                         `json:"metaData"`
}

// ThingEventMsg 上报事件的数据结构
type ThingEventMsg struct {
	ID       string     `json:"id"`
	Version  string     `json:"version"`
	Params   *EventData `json:"params"`
	Type     string     `json:"type"`
	MetaData MetaData   `json:"metaData"`
}

// EventData 事件参数及time
type EventData struct {
	Value PropertyKV `json:"value"`
	Time  int64      `json:"time"`
}

type Reply struct {
	Code int         `json:"code"`
	ID   string      `json:"id"`
	Data interface{} `json:"data"`
}

type Message struct {
	ID      string     `json:"id"`
	Version string     `json:"version"`
	Type    string     `json:"type"`
	Params  PropertyKV `json:"params"`
}

// Client
type Client interface {
	// Connect 设备连接物联网平台
	Connect() error

	// DisConnect 设备取消连接物联网平台
	DisConnect()

	// PubProperty 推送设备属性
	PubProperty(ctx context.Context, meta PropertyKV) (*Reply, error)

	// PubEvent 推送设备事件
	PubEvent(ctx context.Context, meta PropertyKV, eventIdentifier string) (*Reply, error)

	// SubDeviceControl 订阅 topic，获取下行数据，对设备进行调节
	SubDeviceControl(serviceIdentifier string)

	// UnSubDeviceControl 取消订阅
	UnSubDeviceControl(serviceIdentifier string) error

	// normal publish
	// Publish(topic string, data []byte) (*Reply, error)
}
