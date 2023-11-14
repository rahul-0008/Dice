package core

import (
	"errors"
	"io"
	"log"
)

func evalPing(args []string, c io.ReadWriter) error {
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

func EvaluateandRespond(cmd *Rediscmd, c io.ReadWriter) error {
	log.Println(cmd.Cmd)

	switch cmd.Cmd {
	case "PING":
		return evalPing(cmd.Args, c)

	default:
		return evalPing(cmd.Args, c)

	}

}
