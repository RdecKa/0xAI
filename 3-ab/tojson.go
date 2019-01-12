package ab

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
)

// MarshalJSON implements Marshaler interface
// It returns NodeValue in JSON format
func (anv NodeValue) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("{")
	jsonValue, err := json.Marshal(anv.state)
	if err != nil {
		return nil, err
	}
	var val string
	if math.IsInf(anv.value, 1) {
		val = "\"inf\""
	} else if math.IsInf(anv.value, -1) {
		val = "\"-inf\""
	} else {
		val = fmt.Sprintf("%f", anv.value)
	}
	buffer.WriteString(fmt.Sprintf("\"state\":%s,\"val\":%s,\"com\":\"%s\"", jsonValue, val, anv.comment))
	buffer.WriteString("}")
	return buffer.Bytes(), nil
}
