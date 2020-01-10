package index

import "context"

/**
* @Author: hexing
* @Date: 19-9-9 上午11:32
 */
type DownReply func(*Reply)
type SetProperty func(meta Metadata) (Metadata, error)
type ServiceHandle func(name string, meta Metadata) (Metadata, error)
type Options struct {
	//DeviceId      string
	Token         string
	Server        string
	SetProperty   SetProperty
	ServiceHandle ServiceHandle
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
	//device
	PubPropertySync(ctx context.Context, meta Metadata) (*Reply, error) //post property sync
	PubPropertyAsync(meta Metadata) (ReplyChan, error)                  //post property async
	//PubPropertyAsyncEx(meta Metadata, t int64) (ReplyChan, error)                  //post property async
	PubEventSync(ctx context.Context, event string, meta Metadata) (*Reply, error) //post property　sync
	PubEventAsync(event string, meta Metadata) (ReplyChan, error)                  //post property　async
	//driver
	PubTopicPropertySync(ctx context.Context, deviceId, thingId string, meta Metadata) (*Reply, error)            //post property sync
	PubTopicEventSync(ctx context.Context, deviceId, thingId string, event string, meta Metadata) (*Reply, error) //post property　sync
}
