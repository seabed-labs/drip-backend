package utils

import "time"

func GetBoolPtr(val bool) *bool {
	temp := val
	return &temp
}

func GetStringPtr(val string) *string {
	temp := val
	return &temp
}

func GetTimePtr(val time.Time) *time.Time {
	temp := val
	return &temp
}

func GetIntPtr(val int32) *int32 {
	temp := val
	return &temp
}
