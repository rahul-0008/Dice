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

func evict() {
	switch config.EvictingStrtegy {
	case "simple-first":
		evictFirst()
	}
}
