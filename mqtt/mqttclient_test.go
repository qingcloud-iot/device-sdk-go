package mqtt

import (
	"fmt"
	"git.internal.yunify.com/tools/device-sdk-go/index"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"os"
	"testing"
	"time"
)

/**
* @Author: hexing
* @Date: 19-9-9 下午12:19
 */
func TestNewMqtt(t *testing.T) {
	options := &index.Options{
		Token:    "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY3IiOiIxIiwiYXVkIjoiaWFtIiwiYXpwIjoiaWFtIiwiY3VpZCI6ImlhbXItY2x6Y3JrZ2giLCJlaXNrIjoiZXdVX20tdkw4bUxNZnNaQkFDVThtc1VmQmJudDhOVS1EeURwSmZ1Z0Ezdz0iLCJleHAiOjE2MDA0MTMzMDQsImlhdCI6MTU2ODg3NzMwNCwiaXNzIjoic3RzIiwianRpIjoiQTVMVjAzT1Bsc1ZuYndZa1R4Z2ZrYSIsIm5iZiI6MCwib3JnaSI6ImlvdGQtZTI0NmQzMTEtMmIyNC00OTRjLTg1YmUtYTA0Njk2Y2Q3NmMzIiwib3d1ciI6InVzci1rTFZWQkRxZCIsInByZWYiOiJxcm46cWluZ2Nsb3VkOmlhbToiLCJydHlwIjoicm9sZSIsInN1YiI6InN0cyIsInRoaWQiOiJpb3R0LWIwYjQ2YjZiLWNmZDQtNGM3Ny05NTRiLTkzNzU0ZDY3NTUyNCIsInR5cCI6IklEIn0.LwO1mF8iRtdci2QNF3PuqWHSIOzKHOtpcEzecVA8C8kkGdd3bKrzt9b6DvxxfMLh7iEZMl7cwE2vpTRaC5hKJNUyoTiOPr0bAUmcQJngQDAvnR3UC8cY2_AGWiC6tq4778CZ1F2elytgxpDG3oJi85HCMuRyDW0kaCIER2vfY3elPsdmji4EyVeU5sOVJezrVzucvtNI1-_DrQew0MUU3XnT8JkY3px_Nkv6j9CtN3nnR7X18uO8hcUAF0GdzXWcKRDK46b4ZdPSrlF_74umGDH0iLBISRHIACj783jhKqqQqH73jAJfcaWqdUdyZkTp-hHnX7k07gr1DQgS6wNv0Q",
		DeviceId: "iotd-e246d311-2b24-494c-85be-a04696cd76c3",
		Server:   "tcp://192.168.14.120:1883",
		SetProperty: func(meta index.Metadata) (index.Metadata, error) {
			fmt.Println("SetProperty", meta)
			data := make(index.Metadata)
			return data, nil
		},
		ServiceHandle: func(name string, meta index.Metadata) (index.Metadata, error) {
			fmt.Println("ServiceHandle", name, meta)
			data := make(index.Metadata)
			return data, nil
		},
	}
	m, err := NewMqtt(options)
	assert.Nil(t, err)
	assert.NotNil(t, m)
	time.Sleep(5 * time.Second)

	go func() {
		data := index.Metadata{
			"int32":  10,
			"float":  rand.Float32(),
			"double": rand.Float64(),
			"string": "xxxxxxxxxxxxxxxxx",
		}
		ch, err := m.PubPropertyAsync(data)
		//ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		//reply, err := m.PubPropertySync(ctx, data)
		//cancel()
		assert.Nil(t, err)
		select {
		case value := <-ch:
			fmt.Println(value)
		}
		//data = index.Metadata{
		//	"int32":  10,
		//	"string": "hexing-string",
		//	"float":  rand.Float32(),
		//	"double": rand.Float64(),
		//}
		//reply, err = m.PubEventSync(context.Background(), "he-event1", data)
		//assert.Nil(t, err)
		//fmt.Println(reply)
		time.Sleep(5 * time.Second)
		os.Exit(0)
	}()
	select {}
	//name := "test"
	//reply = m.PubEvent(ctx, name, data)
	//t.Log(reply)
}
