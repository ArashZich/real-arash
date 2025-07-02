package endpoints

import (
	"errors"
	"math"
	"math/rand"
	"strings"
	"time"
)

const (
	Numbers      = "0123456789"
	Alphabetic   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	Alphanumeric = Numbers + Alphabetic
	MinCount     = 1
	MinLength    = 6
	PatternChar  = "#"
)

func randomInt(min, max int) int {
	return min + rand.Intn(1+max-min)
}

func randomChar(cs []byte) string {
	return string(cs[randomInt(0, len(cs)-1)])
}

func repeatStr(count int, str string) string {
	return strings.Repeat(str, count)
}

func numberOfChar(str, char string) int {
	return strings.Count(str, char)
}

func isFeasible(charset, pattern, char string, count int) bool {
	ls := numberOfChar(pattern, char)
	return math.Pow(float64(len(charset)), float64(ls)) >= float64(count)
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func generateCode(length, count int, charset, pattern, prefix, suffix string) (string, error) {

	if length == 0 {
		length = numberOfChar(pattern, PatternChar)
	}
	if count == 0 {
		return "", errors.New("تعداد نمیتواند صفر باشد")
	}
	if len(charset) == 0 {
		return "", errors.New("مجموعه کارکترها نمی تواند خالی باشد")
	}
	if pattern == "" {
		return "", errors.New("الگو نمیتواند خالی باشد")
	}

	numPatternChar := numberOfChar(pattern, PatternChar)
	if length == 0 || length != numPatternChar {
		length = numPatternChar
	}

	codes := make([]string, count)
	for i := 0; i < count; i++ {
		pts := strings.Split(pattern, "")
		for i, v := range pts {
			if v == PatternChar {
				pts[i] = randomChar([]byte(charset))
			}
		}
		codes[i] = prefix + strings.Join(pts, "") + suffix
	}

	return strings.Join(codes, ""), nil
}
