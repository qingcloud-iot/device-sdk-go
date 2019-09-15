package mqtt

import (
	"git.internal.yunify.com/tools/device-sdk-go/index"
	"github.com/stretchr/testify/assert"
	"testing"
)

/**
* @Author: hexing
* @Date: 19-9-9 下午12:19
 */
func TestNewMqtt(t *testing.T) {
	options := &index.Options{}
	m, err := NewMqtt(options)
	assert.Nil(t, err)
	assert.NotNil(t, m)
	err = m.Start(nil, nil)
	assert.Nil(t, err)
	select {}
}
