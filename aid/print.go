package aid

import (
	"encoding/json"
	"fmt"
)

func PrintJSON(v interface{}) {
	json1, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(json1))
}