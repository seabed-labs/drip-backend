package utils

/*
GetWithDefault - Return value if its not nil, else return a fallback
ex:
```
const value := "my value"
var nilValue *string = nil
const fallback = "my fallback"

assert.Equal(t, GetWithDefault(&value, fallback), value)
assert.Equal(t, GetWithDefault(nilValue, fallback), fallback)
```
*/
func GetWithDefault[A any](value *A, fallBack A) A {
	if value != nil {
		return *value
	}
	return fallBack
}
