package mqtt

/**
* @Author: hexing
* @Date: 19-9-15 下午3:24
 */
import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"strings"
)

const (
	TOKEN_DEVICE_ID = "orgi"
	TOKEN_THING_ID  = "thid"
	WORKER_ID       = 1
	MQTT_VERSION    = "v1.0.0"
	WORKER_POOL     = 10
)

// return token payload
func parseToken(deviceToken string) (string, string, error) {
	var (
		deviceId string
		thingId  string
		err      error
	)
	defer func() {
		if errToken := recover(); err != nil {
			err = errToken.(error)
			return
		}
	}()
	token, _ := jwt.Parse(deviceToken, func(token *jwt.Token) (interface{}, error) {
		return "", nil
	})
	if payload, ok := token.Claims.(jwt.MapClaims); ok {
		if err := payload.Valid(); err != nil {
			return "", "", err
		}
		deviceId, ok = payload[TOKEN_DEVICE_ID].(string)
		if !ok {
			return "", "", errors.New("device id type error")
		}
		thingId, ok = payload[TOKEN_THING_ID].(string)
		if !ok {
			return "", "", errors.New("device id type error")
		}
		return deviceId, thingId, nil
	} else {
		return deviceId, thingId, errors.New("token error")
	}
}
func parseServiceName(topic string) string {
	kv := strings.Split(topic, "/")
	if len(kv) != 8 {
		return ""
	}
	return kv[6]
}

func isServiceTopic(thingId, deviceId, topic string) bool {
	temp := fmt.Sprintf("/sys/%s/%s/thing/service", thingId, deviceId)
	return strings.HasPrefix(topic, temp)
}
