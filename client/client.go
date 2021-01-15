package client

import (
	"context"
	"github.com/qingcloud-iot/device-sdk-go/define"
)

// Client 设备 sdk
type Client interface {

	// Connect 设备连接物联网平台
	Connect() error

	// DisConnect 设备取消连接物联网平台
	DisConnect()

	// PubProperty 推送设备属性
	PubProperty(ctx context.Context, meta define.PropertyKV) (*define.Reply, error)

	// PubPropertyWithTime
	PubPropertyWithTime(ctx context.Context, metaWithTime define.PropertyKVWithTime) (*define.Reply, error)

	// PubEvent 推送设备事件
	PubEvent(ctx context.Context, meta define.PropertyKV, eventIdentifier string) (*define.Reply, error)

	// SubDeviceControl 订阅 topic，获取下行数据，对设备进行调节
	SubDeviceControl(serviceIdentifier string)

	// UnSubDeviceControl 取消订阅
	UnSubDeviceControl(serviceIdentifier string) error

	// normal subscribe
	Subscribe(topic string, qos int32, cb MessageCallback) error

	SubscribeMultiple(topics []string, cb MessageCallback) error

	Unsubscribe(topics []string) error

	Publish(topic string, qos int32, payload []byte) error
}
