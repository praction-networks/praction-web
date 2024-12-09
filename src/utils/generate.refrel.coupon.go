package utils

import (
	"math/rand"
	"time"
)

const (
	couponLength = 9
	upperLetters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digits       = "0123456789"
	characters   = upperLetters + digits
)

// GenerateRefrelCoupon generates a random 9-character coupon consisting of uppercase letters and digits.
func GenerateRefrelCoupon() string {
	// Use NewSource and New to create a new rand.Rand generator
	source := rand.NewSource(time.Now().UnixNano()) // Generate a new source with current time
	r := rand.New(source)                           // Create a new random number generator with the source

	coupon := make([]byte, couponLength)

	for i := range coupon {
		coupon[i] = characters[r.Intn(len(characters))] // Use r.Intn to get a random character
	}

	return string(coupon)
}
