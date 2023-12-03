package core

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/DiceDB/Dice/config"
)

// TODO: Support Expiration
// TODO: Support non-kv data structures
// TODO: Support sync write
func dumpKey(fp *os.File, k string, obj *Obj) {
	cmd := fmt.Sprintf("SET %s %s", k, obj.Value)
	tokens := strings.Split(cmd, " ")
	fp.Write(Encode(tokens, false))
}

// TODO: To da new and switch
func DumpAllAOF() {
	fp, err := os.OpenFile(config.AOFFile, os.O_CREATE|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Println("error", err)
		return
	}
	for k, obj := range store {
		dumpKey(fp, k, obj)
	}
	log.Println("AOF file write complete")
}
