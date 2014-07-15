package signaling

import "encoding/json"
import "bytes"

func ToJsonString(info *map[string]string) string {
	var buf bytes.Buffer
	result, _ := json.Marshal(info)
	buf.Write(result)
	return buf.String()
}
