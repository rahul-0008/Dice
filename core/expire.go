package core

import (
	"time"
)

func expireSample() float32 {
	var limit = 20
	var expiredCount = 0

	// assusmong the go lang hash is randomized
	for key, obj := range store {
		if obj.ExpiresAt != -1 {
			limit--

			if obj.ExpiresAt <= time.Now().UnixMilli() {
				Del(key)
				expiredCount++
			}
		}
		if limit == 0 {
			break
		}
	}

	return float32(expiredCount) / float32(20.0)

}

// Deletes all the expired keys - the active way
// Sampling approach: https://redis.io/commands/expire/
func DeleteExpiredKeys() {
	for {
		frac := expireSample()

		if frac < 0.25 {
			break
		}

	}

}
