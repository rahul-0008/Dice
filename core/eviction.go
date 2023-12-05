package core

import (
	"github.com/DiceDB/Dice/config"
)

func evictFirst() {

	for k := range store {
		delete(store, k)
		break
	}
}

// Randomly removes keys to make space for the new data added.
// The number of keys removed will be sufficient to free up least 10% space
func evictAllKeysRandom() {
	evictCount := int64(config.EvictionRatio * float64(config.KeysLimit))
	// Iteration of Golang dictionary can be considered as a random
	// because it depends on the hash of the inserted key
	for k := range store {
		Del(k)
		evictCount--
		if evictCount <= 0 {
			break
		}
	}
}

func evict() {
	switch config.EvictingStrtegy {
	case "simple-first":
		evictFirst()
	case "allkeys-random":
		evictAllKeysRandom()
	}
}
