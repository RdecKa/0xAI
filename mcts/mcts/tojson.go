package mcts

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// MarshalJSON implements Marshaler interface
// It returns MCTS in JSON format
func (mcts MCTS) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("{")
	jsonValue, err := json.Marshal(mcts.mcTree)
	if err != nil {
		return nil, err
	}
	buffer.WriteString(fmt.Sprintf("\"tree\":%s, \"c\":%f", jsonValue, mcts.c))
	buffer.WriteString("}")
	return buffer.Bytes(), nil
}

// MarshalJSON implements Marshaler interface
// It returns mctsNodeValue in JSON format
func (mnv mctsNodeValue) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("{")
	jsonValue, err := json.Marshal(mnv.state)
	if err != nil {
		return nil, err
	}
	buffer.WriteString(fmt.Sprintf("\"N\": %d, \"Q\": %f, \"state\": %s", mnv.n, mnv.q, jsonValue))
	buffer.WriteString("}")
	return buffer.Bytes(), nil
}

// WriteToFile saves MCTS in JSON format to the file folder/currentDate
func WriteToFile(mcts MCTS, folder string, filePrefix string, addIndent bool) error {
	// Create a new file
	t := time.Now()
	fileName := folder + "/" + filePrefix + "_" + t.Format("2006-01-02T15:04:05") + ".json"
	fmt.Printf("Writing MCTS to file %s ... ", fileName)
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	// Create JSON
	jsonText, err := json.Marshal(mcts)
	if err != nil {
		return err
	}

	// Add indent
	if addIndent {
		var out bytes.Buffer
		json.Indent(&out, jsonText, "", "\t")
		jsonText = out.Bytes()
	}

	// Write to file
	_, err = f.Write(jsonText)
	if err != nil {
		return err
	}
	fmt.Println("Done!")
	return nil
}
