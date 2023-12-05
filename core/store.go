package core

import (
	"log"
	"time"

	"github.com/DiceDB/Dice/config"
)

var store map[string]*Obj

func init() {
	store = make(map[string]*Obj)
}

func NewObj(value interface{}, durationMs int64, oType uint8, oEnc uint8) *Obj {
	log.Println("Time to live given ", durationMs)
	var expiresAt int64 = -1
	if durationMs > 0 {
		expiresAt = time.Now().UnixMilli() + durationMs
	}

	return &Obj{
		Value:        value,
		ExpiresAt:    expiresAt,
		TypeEncoding: oType | oEnc,
	}
}

func Put(K string, obj *Obj) {
	if len(store) > config.KeysLimit {
		evict()
	}
	store[K] = obj
	if KeyspaceStat[0] == nil {
		KeyspaceStat[0] = make(map[string]int)
	}
	KeyspaceStat[0]["keys"]++
}

func Get(K string) *Obj {
	v := store[K]
	if v != nil {
		if v.ExpiresAt != -1 && v.ExpiresAt <= time.Now().UnixMilli() {
			Del(K)
			return nil
		}
	}
	return v
}

func Del(K string) bool {
	if _, ok := store[K]; ok {
		delete(store, K)
		KeyspaceStat[0]["keys"]--
		return true
	}
	return false

}
