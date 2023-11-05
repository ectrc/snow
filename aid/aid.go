package aid

import (
	"encoding/json"
	"os"
	"os/signal"
	"syscall"
)

func WaitForExit() {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
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