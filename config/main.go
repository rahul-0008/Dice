package config

var Host string = "0.0.0.0"
var Port int = 7379
var KeysLimit = 5

var EvictingStrtegy string = "simple-first"
var AOFFile string = "./dice-master.aof"
