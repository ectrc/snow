package aid

import "github.com/goccy/go-json"

type JSON map[string]interface{}

func JSONFromBytes(input []byte) JSON {
	var output JSON
	json.Unmarshal(input, &output)
	return output
}

func (j *JSON) ToBytes() []byte {
	json, _ := json.Marshal(j)
	return json
}

func JSONToBytes(input JSON) []byte {
	json, _ := json.Marshal(input)
	return json
}

func JSONStringify(input interface{}) string {
	json, _ := json.Marshal(input)
	return string(json)
}

func JSONParse(input string) interface{} {
	var output interface{}
	json.Unmarshal([]byte(input), &output)
	return output
}