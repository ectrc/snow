package aid

import "math/rand"

var Random *rand.Rand

func SetRandom(r *rand.Rand) {
	Random = r
}

func RandomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	s := make([]rune, n)
	for i := range s {
 		s[i] = letters[Random.Intn(len(letters))]
	}

	return string(s)
}

func RandomInt(min, max int) int {
	return Random.Intn(max-min) + min
}
