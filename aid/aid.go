package aid

import (
	m "math/rand"
	"os"
	"os/signal"
	"regexp"
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

func RandomInt(min, max int) int {
	return m.Intn(max - min) + min
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

func Regex(str, regex string) *string {
	// reg := regexp.MustCompile(`(?:CID_)(\d+|A_\d+)(?:_.+)`).FindStringSubmatch(strings.Join(split[:], "_"))
	reg := regexp.MustCompile(regex).FindStringSubmatch(str)
	if len(reg) > 1 {
		return &reg[1]
	}

	return nil
}

func Ternary[T any](condition bool, a, b T) T {
	if condition {
		return a
	}
	return b
}