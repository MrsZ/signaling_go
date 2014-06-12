package signaling

import "log"
import "encoding/json"
import "io"
import "bytes"


func ReadJson(from io.Reader, to interface{}) error {
	dec := json.NewDecoder(from)
	if err := dec.Decode(to); err != nil {
		log.Printf("Failed to parse json: %s", err)
		return err
	}
	return nil
}

func ToJsonString(info *map[string]string) string {
	var buf bytes.Buffer
	result, _ := json.Marshal(info)
	buf.Write(result)
	return buf.String()
}
