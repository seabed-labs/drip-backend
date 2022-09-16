package controller

func hasValue(params []string, value string) bool {
	for _, v := range params {
		if v == value {
			return true
		}
	}
	return false
}
