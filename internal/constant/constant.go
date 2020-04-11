package constant

const (
	MQTT_VERSION = "1.0"
)

const (
	// TOKEN_DEVICE_ID token payload 的 entityID 字段
	TOKEN_DEVICE_ID = "orgi"
	// TOKEN_THING_ID token payload 的 modelID 字段
	TOKEN_THING_ID = "thid"
)

const (
	// POST_PROPERTY_TOPIC 属性上报的 topic，数据上行
	POST_PROPERTY_TOPIC = "/sys/%s/%s/thing/property/%s/post"

	// POST_EVENT_TOPIC 事件上报的 topic，数据上行
	POST_EVENT_TOPIC = "/sys/%s/%s/thing/event/%s/post"

	// SET_PROPERTY_TOPIC_REPLY 属性上报成功后的 reply
	SET_PROPERTY_TOPIC_REPLY = "/sys/%s/%s/thing/service/set/call_reply"

	// DEVICE_CONTROL_TOPIC 服务调用订阅的 topig，数据下行
	DEVICE_CONTROL_TOPIC = "/sys/%s/%s/thing/service/%s/call"

	// SET_SERVICE_TOPIC_REPLY 服务调用成功后的 reply
	SET_SERVICE_TOPIC_REPLY = "/sys/%s/%s/thing/service/%s/call_reply"
)

const (
	RPC_SUCCESS = 200
	RPC_TIMEOUT = 1001
)

const (
	// PROPERTY_TYPE_BASE 上报基础属性时，用于构造 topic 的常量
	PROPERTY_TYPE_BASE = "base"

	// PROPERTY_TYPE 属性上报时的 type 字段值
	PROPERTY_TYPE = "thing.property.post"
)
