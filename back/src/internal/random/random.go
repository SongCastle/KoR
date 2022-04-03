package random

import (
	"crypto/rand"
	"encoding/base64"
	mrand "math/rand"
	"strconv"
	"time"
)

var Encoder = base64.RawStdEncoding

func Generate(len int) string {
	if r, err := generateStrongly(len); err == nil {
		return r
	}
	return generateWeakly(len)
}

func generateStrongly(len int) (string, error) {
	b := make([]byte, len)
	_, err := rand.Read(b)
	if err == nil {
		return Encoder.EncodeToString(b)[:len], nil
	}
	return "", err
}

func generateWeakly(len int) string {
	mrand.Seed(time.Now().UnixNano())

	word := ""
	for i := 0; i < len; i++ {
		n := mrand.Intn(62)
		if n < 26 {
			// 小文字
			word += string('a' + n)
		} else if n < 52 {
			// 大文字
			word += string('A' + (n - 26))
		} else {
			// 数字
			word += strconv.Itoa(n - 52)
		}
	}
	return word
}
