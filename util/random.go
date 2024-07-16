package util

import (
	"math/rand"
	"time"
)

// init initializes the random number generator.
func init() {
	rand.NewSource(time.Now().UnixNano())
}

// RandomString generates a random string of given length.
func RandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

// RandomDate generates a random date within a specified range.
func RandomDate(start, end time.Time) time.Time {
	delta := end.Sub(start)
	sec := rand.Int63n(delta.Nanoseconds() / 1000000000)
	return start.Add(time.Duration(sec) * time.Second)
}

// RandomBool generates a random boolean value.
func RandomBool() bool {
	return rand.Intn(2) == 0
}
