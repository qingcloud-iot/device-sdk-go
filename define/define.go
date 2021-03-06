package define

/*
	定义消息结构对应的 go 类型
*/

// PropertyKV 事件 kv 健值对
type PropertyKV map[string]interface{}

// PropertyKVWithTime 允许用户自定义时间
type PropertyKVWithTime map[string]*PropertyValueAndTime

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
	MetaData MetaData                         `json:"metadata"`
}

// ThingEventMsg 上报事件的数据结构
type ThingEventMsg struct {
	ID       string     `json:"id"`
	Version  string     `json:"version"`
	Params   *EventData `json:"params"`
	Type     string     `json:"type"`
	MetaData MetaData   `json:"metadata"`
}

// EventData 事件参数及time
type EventData struct {
	Value PropertyKV `json:"value"`
	Time  int64      `json:"time"`
}

// Message 数据下行，下发的数据结构
type Message struct {
	ID      string     `json:"id"`
	Version string     `json:"version"`
	Type    string     `json:"type"`
	Params  PropertyKV `json:"params"`
}

// Reply client 返回的数据结构
type Reply struct {
	Code int         `json:"code"`
	ID   string      `json:"id"`
	Data interface{} `json:"data"`
}
