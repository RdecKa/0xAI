package hex

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// MarshalJSON implements Marshaler interface
// It returns State in JSON format
func (state State) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("{")
	jsonValue, err := json.Marshal(state.grid)
	if err != nil {
		return nil, err
	}
	buffer.WriteString(fmt.Sprintf("\"grid\":%s,\"lp\":\"%s\"", jsonValue, state.lastAction.c))
	buffer.WriteString("}")
	return buffer.Bytes(), nil
}

// MarshalJSON implements Marshaler interface
// It returns Action in JSON format
func (action Action) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("{")
	buffer.WriteString(fmt.Sprintf("\"x\":%d,\"y\":%d,\"c\":\"%v\"", action.x, action.y, action.c))
	buffer.WriteString("}")
	return buffer.Bytes(), nil
}
