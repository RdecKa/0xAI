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
	if string(jsonValue) == "{}" {
		buffer.WriteString("}")
		return buffer.Bytes(), nil
	}
	buffer.WriteString(fmt.Sprintf("\"value\": %s,", jsonValue))

	buffer.WriteString("\"children\":[")
	firstWritten := false
	for _, c := range node.children {
		jsonValue, err = json.Marshal(c)
		if err != nil {
			return nil, err
		}
		if string(jsonValue) != "{}" {
			if !firstWritten {
				firstWritten = true
			} else {
				buffer.WriteString(",")
			}
			buffer.WriteString(fmt.Sprintf("%s", jsonValue))
		}
	}
	buffer.WriteString("]")

	buffer.WriteString("}")
	return buffer.Bytes(), nil
}
