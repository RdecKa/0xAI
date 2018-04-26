package tree

import (
	"bytes"
	"fmt"
)

// MarshalJSON implements Marshaler interface
// It returns Tree in JSON format (TODO)
func (tree Tree) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("{")
	buffer.WriteString(fmt.Sprintf("\"root\":\"%s\"", "fdg"))
	buffer.WriteString("}")
	return buffer.Bytes(), nil
}
