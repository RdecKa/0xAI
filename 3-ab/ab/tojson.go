package ab

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// MarshalJSON implements Marshaler interface
// It returns AbNodeValue in JSON format
func (anv AbNodeValue) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("{")
	jsonValue, err := json.Marshal(anv.state)
	if err != nil {
		return nil, err
	}
	buffer.WriteString(fmt.Sprintf("\"tree\":%s, \"val\":%f", jsonValue, anv.value))
	buffer.WriteString("}")
	return buffer.Bytes(), nil
}
