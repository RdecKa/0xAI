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
	buffer.WriteString(fmt.Sprintf("\"state\":%s,\"val\":%f,\"comment\":\"%s\"", jsonValue, anv.value, anv.comment))
	buffer.WriteString("}")
	return buffer.Bytes(), nil
}

// MarshalJSON implements Marshaler interface
// It returns RootNodeValue in JSON format
/*func (rnv RootNodeValue) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("{")
	jsonValueRoot, err := json.Marshal(rnv.root)
	if err != nil {
		return nil, err
	}
	jsonValueState, err := json.Marshal(rnv.state)
	if err != nil {
		return nil, err
	}
	buffer.WriteString(fmt.Sprintf("\"rootAB\":%s,\"state\":%s,\"size\":%d", jsonValueRoot, jsonValueState, rnv.size))
	buffer.WriteString("}")
	return buffer.Bytes(), nil
}*/
