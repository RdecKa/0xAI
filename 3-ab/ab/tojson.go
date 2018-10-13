package ab

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// MarshalJSON implements Marshaler interface
// It returns NodeValue in JSON format
func (anv NodeValue) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("{")
	jsonValue, err := json.Marshal(anv.state)
	if err != nil {
		return nil, err
	}
	buffer.WriteString(fmt.Sprintf("\"state\":%s,\"val\":%f,\"com\":\"%s\"", jsonValue, anv.value, anv.comment))
	buffer.WriteString("}")
	return buffer.Bytes(), nil
}
