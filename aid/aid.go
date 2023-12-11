package aid

import (
	"math/rand"
	"os"
	"os/signal"
	"syscall"

	"github.com/goccy/go-json"
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

func RandomString(n int) string {
    var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
 
    s := make([]rune, n)
    for i := range s {
        s[i] = letters[rand.Intn(len(letters))]
    }
    return string(s)
}