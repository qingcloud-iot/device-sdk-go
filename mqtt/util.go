package mqtt

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"math/rand"
	"strings"
)

const (
	TOKEN_DEVICE_ID = "orgi"
	TOKEN_THING_ID  = "thid"
	WORKER_ID       = 1
	MQTT_VERSION    = "1.0"
	WORKER_POOL     = 10
)

// return token payload
func parseToken(deviceToken string) (string, string, error) {
	var (
		entityID string
		modelID  string
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
		entityID, ok = payload[TOKEN_DEVICE_ID].(string)
		if !ok {
			return "", "", errors.New("device id type error")
		}
		modelID, ok = payload[TOKEN_THING_ID].(string)
		if !ok {
			return "", "", errors.New("device id type error")
		}
		return entityID, modelID, nil
	} else {
		return entityID, modelID, errors.New("token error")
	}
}
func parseServiceName(topic string) string {
	kv := strings.Split(topic, "/")
	if len(kv) != 8 {
		return ""
	}
	return kv[6]
}

func isServiceTopic(modelID, entityID, topic string) bool {
	temp := fmt.Sprintf("/sys/%s/%s/thing/service", modelID, entityID)
	return strings.HasPrefix(topic, temp)
}
func RandInt64(min, max int64) int64 {
	if min >= max || min == 0 || max == 0 {
		return max
	}
	return rand.Int63n(max-min) + min
}
