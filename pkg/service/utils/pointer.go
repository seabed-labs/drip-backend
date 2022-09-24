package utils

func GetBoolPtr(val bool) *bool {
	temp := val
	return &temp
}

func GetStringPtr(val string) *string {
	temp := val
	return &temp
}
