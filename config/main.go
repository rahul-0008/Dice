package config

var Host string = "0.0.0.0"
var Port int = 7379
var KeysLimit = 100

// will evict EvictRatio of keys whenever eviction runs
var EvictionRatio = 0.40

var EvictingStrtegy string = "allkeys-random"
var AOFFile string = "./dice-master.aof"
