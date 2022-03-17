package random

import (
	// TODO: "crypto/rand" にする
	"math/rand"
	"time"
	"strconv"
)

func Generate(len int) string {
	rand.Seed(time.Now().UnixNano())

	word := ""
	for i := 0; i < len; i++ {
		n := rand.Intn(62)
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
