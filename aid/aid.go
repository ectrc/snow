package aid

import (
	m "math/rand"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

func WaitForExit() {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}

func RandomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	s := make([]rune, n)
	for i := range s {
		s[i] = letters[m.Intn(len(letters))]
	}
	return string(s)
}

func FormatNumber(number int) string {
	str := ""
	for i, char := range ReverseString(strconv.Itoa(number)) {
		if i % 3 == 0 && i != 0 {
			str += ","
		}
		str += string(char)
	}

	return ReverseString(str)
}

func ReverseString(input string) string {
	str := ""
	for _, char := range input {
		str = string(char) + str
	}
	return str
}

func ToHex(number int) string {
	inta := strconv.FormatInt(int64(number), 16)
	
	if len(inta) == 1 {
		return "0" + inta
	}
	
	return inta
}