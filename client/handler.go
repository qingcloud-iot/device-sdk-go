package client

import (
	"github.com/qingcloud-iot/device-sdk-go/define"
)

// CallBack 回调
type CallBack interface {
	// Handler 处理下发的消息
	Handler(msg *define.Message) define.PropertyKV
}

// MessageCallback 消息回调
type MessageCallback func(string, []byte)
