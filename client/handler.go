package client

import (
	"git.internal.yunify.com/iot-sdk/device-sdk-go/define"
)

// CallBack 回调
type CallBack interface {
	// Handler 处理下发的消息
	Handler(msg *define.Message) define.PropertyKV
}
