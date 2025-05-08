package common

import "errors"

type DelayType string

const (
	DelayTypeHour   DelayType = "hour"
	DelayTypeMinute DelayType = "minute"
)

var ErrDelayType = errors.New("incorrect delay type")

func GetDelayType(value string) (DelayType, error) {
	switch value {
	case string(DelayTypeHour):
		return DelayTypeHour, nil
	case string(DelayTypeMinute):
		return DelayTypeMinute, nil
	default:
		return "", ErrDelayType
	}
}
