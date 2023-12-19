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