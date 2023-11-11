package core

import (
	"errors"
	"log"
	"net"
)

func evalPing(args []string, c net.Conn) error {
	var b []byte

	if len(args) >= 2 {
		return errors.New("ERR wrong number of arguments for 'ping' command")
	}

	if len(args) == 0 {
		simpleString := true
		b = Encode("PONG", simpleString)
	} else {
		simpleString := false
		b = Encode(args[0], simpleString)
	}

	_, err := c.Write(b)
	return err
}

func EvaluateandRespond(cmd *Rediscmd, c net.Conn) error {
	log.Println(cmd.Cmd)

	switch cmd.Cmd {
	case "PING":
		return evalPing(cmd.Args, c)

	default:
		return evalPing(cmd.Args, c)

	}

}
