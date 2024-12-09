package utils

import (
	"math/rand"
	"time"
)

func GenerateRandomOTP(length int) int64 {
	// Create a new random generator with a unique seed
	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)

	// Calculate the minimum and maximum range for the OTP
	min := int64(1)
	for i := 1; i < length; i++ {
		min *= 10
	}
	max := min*10 - 1

	// Generate a random number in the range [min, max]
	return random.Int63n(max-min+1) + min
}
