package client

import (
	"context"
	"git.internal.yunify.com/iot-sdk/device-sdk-go/define"
)

// Client 设备 sdk
type Client interface {

	// Connect 设备连接物联网平台
	Connect() error

	// DisConnect 设备取消连接物联网平台
	DisConnect()

	// PubProperty 推送设备属性
	PubProperty(ctx context.Context, meta define.PropertyKV) (*define.Reply, error)

	// PubEvent 推送设备事件
	PubEvent(ctx context.Context, meta define.PropertyKV, eventIdentifier string) (*define.Reply, error)

	// SubDeviceControl 订阅 topic，获取下行数据，对设备进行调节
	SubDeviceControl(serviceIdentifier string)

	// UnSubDeviceControl 取消订阅
	UnSubDeviceControl(serviceIdentifier string) error

	// normal publish
	// Publish(topic string, data []byte) (*Reply, error)
}
