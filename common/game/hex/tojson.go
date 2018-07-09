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
	buffer.WriteString(fmt.Sprintf("\"grid\": %s, \"lastPlayer\": \"%s\"", jsonValue, state.lastPlayer))
	buffer.WriteString("}")
	return buffer.Bytes(), nil
}
