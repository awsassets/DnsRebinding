package random

import (
	"math/rand"
	"strings"
	"time"
)

const (
	LowerCharset  = "abcdefghijklmnopqrstuvwxyz"
	UpperCharset  = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	NumberCharset = "0123456789"
)

func Int(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min+1) + min
}

func StringWithCharset(length int, charset string) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func String(length int) string {
	return StringWithCharset(length, LowerCharset+UpperCharset+NumberCharset)
}

func Upper(s string) string {
	rand.Seed(time.Now().UnixNano())
	r := ""
	for i := range s {
		if rand.Intn(2) > 0 {
			r += strings.ToUpper(string(s[i]))
			continue
		}
		r += string(s[i])
	}
	return r
}
