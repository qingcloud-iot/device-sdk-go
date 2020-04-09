package define

/*
	定义消息结构对应的 go 类型
*/

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
