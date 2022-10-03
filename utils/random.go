package utils

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// seed for generating random string
const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomInt -> random integer between min and max
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// RandomString -> random string of n length
func RandomString(n int) string {
	var sb strings.Builder
	s := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(s)]
		sb.WriteByte(c)
	}
	return sb.String()
}

// RandonOwner -> generates a random owner name
func RandomOwner() string {
	return RandomString(6)
}

// RandomMoney -> generates a randon balance
func RandomMoney() int64 {
	return RandomInt(0, 10000)
}

// RandomCurrency -> generates a random currency from the list
func RandomCurrency() string {
	currency := []string{INR, USD, EUR, CAD, AUD}
	n := len(currency)
	return currency[rand.Intn(n)]
}

// RandomEmail randomly generated email fro testing
func RandomEmail() string {
	return fmt.Sprintf("%s@email.com", RandomString(6))
}
