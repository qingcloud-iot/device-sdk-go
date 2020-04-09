package constant

const (
	TOKEN_DEVICE_ID = "orgi"
	TOKEN_THING_ID  = "thid"
	MQTT_VERSION    = "1.0"
)

const (
	POST_PROPERTY_TOPIC      = "/sys/%s/%s/thing/property/%s/post"       // 属性上报 topic
	POST_EVENT_TOPIC         = "/sys/%s/%s/thing/event/%s/post"          //down
	SET_PROPERTY_TOPIC_REPLY = "/sys/%s/%s/thing/service/set/call_reply" //down
	SET_SERVICE_TOPIC_REPLY  = "/sys/%s/%s/thing/service/%s/call_reply"

	DEVICE_CONTROL_TOPIC = "/sys/%s/%s/thing/service/%s/call"
)

const (
	RPC_SUCCESS = 200  //success
	RPC_TIMEOUT = 1001 // rpc timeout
)


const (
	// 上报系统属性时，用于构造 topic
	PROPERTY_TYPE_BASE = "base"

	// 属性上报时的 type 字段值
	PROPERTY_TYPE = "thing.property.post"
)

