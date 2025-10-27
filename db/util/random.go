package util

import (
	"math/rand"
	"strings"
	"time"
)

func init() {

	rand.New(rand.NewSource(time.Now().UnixNano()))

}

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func RandomInt(min, max int) int {
	return min + rand.Intn(max-min+1)
}

func RandomFloat(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func RandomStr(length int) string {

	var sb strings.Builder
	for i := 0; i < length; i++ {
		sb.WriteByte(alphabet[rand.Intn(len(alphabet))])
	}
	return sb.String()

}

func RandomMoney() int64 {

	return int64(RandomInt(0, 10000))

}

func RandomEntryMoney() int64 {

	return int64(RandomInt(-10000, 10000))

}

func RandomCurrency() string {
	currencies := []string{"USD", "EUR", "GBP", "JPY"}
	return currencies[rand.Intn(len(currencies))]

}
