package mqtt

/**
* @Author: hexing
* @Date: 19-9-15 下午3:24
 */
import (
	"errors"
)
const (
	TOKEN_DEVICE_ID = "orgi"
	TOKEN_THING_ID  = "thid"
)
// return token payload
func parseToken(deviceToken string) (string, string, error) {
	var (
		deviceId string
		thingId string
		err error
	)
	defer func() {
		if errToken := recover();err != nil {
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