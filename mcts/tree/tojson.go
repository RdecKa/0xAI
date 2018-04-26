package tree

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// MarshalJSON implements Marshaler interface
// It returns Tree in JSON format
func (tree Tree) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("{")
	jsonValue, err := json.Marshal(tree.root)
	if err != nil {
		return nil, err
	}
	buffer.WriteString(fmt.Sprintf("\"root\": %s", jsonValue))
	buffer.WriteString("}")
	return buffer.Bytes(), nil
}

// MarshalJSON implements Marshaler interface
// It returns Node in JSON format
func (node Node) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("{")

	jsonValue, err := json.Marshal(node.value)
	if err != nil {
		return nil, err
	}
	buffer.WriteString(fmt.Sprintf("\"value\": %s,", jsonValue))

	jsonValue, err = json.Marshal(node.children)
	if err != nil {
		return nil, err
	}
	buffer.WriteString(fmt.Sprintf("\"children\": %s", jsonValue))

	buffer.WriteString("}")
	return buffer.Bytes(), nil
}
