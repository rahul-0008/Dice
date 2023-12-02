package core

import (
	"errors"
	"io"
	"log"
	"strconv"
	"time"
)

var RESP_NIL []byte = []byte("$-1\r\n")

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
func evalSet(args []string, c io.ReadWriter) error {
	if len(args) <= 1 {
		return errors.New("(error) ERR wrong number of arguments for 'set' command")
	}
	var key, value string
	var exDurationMs int64 = -1

	key, value = args[0], args[1]

	for i := 2; i < len(args); i++ {
		log.Println(args[i])
		switch args[i] {
		case "EX", "ex":
			i++
			if i == len(args) {
				return errors.New("(error) ERR syntax error")
			}
			exDurationSec, err := strconv.ParseInt(args[3], 10, 64)
			if err != nil {
				return errors.New("(error) ERR value is not an integer or out of range")
			}
			exDurationMs = exDurationSec * 1000
		default:
			return errors.New("(error) ERR syntax error")
		}
	}

	//putting the k and v in th Hash
	log.Println("Time to live given ", exDurationMs)
	Put(key, NewObj(value, int64(exDurationMs)))
	c.Write([]byte("+OK\r\n"))
	return nil
}
func evalGet(args []string, c io.ReadWriter) error {
	if len(args) != 1 {
		return errors.New("(error) ERR wrong number of arguments for 'get' command")
	}

	var key string = args[0]

	// Get the key from the hash table
	obj := Get(key)

	// if key does not exist, return RESP encoded nil
	if obj == nil {
		c.Write(RESP_NIL)
		return nil
	}

	// if key already expired then return nil
	if obj.ExpiresAt != -1 && obj.ExpiresAt <= time.Now().UnixMilli() {
		c.Write(RESP_NIL)
		return nil
	}

	// return the RESP encoded value
	c.Write(Encode(obj.Value, false))
	return nil
}

func evalTTL(args []string, c io.ReadWriter) error {
	if len(args) != 1 {
		return errors.New("(error) ERR wrong number of arguments for 'ttl' command")
	}

	var key string = args[0]

	obj := Get(key)

	// if key does not exist, return RESP encoded -2 denoting key does not exist
	if obj == nil {
		c.Write([]byte(":-2\r\n"))
		return nil
	}

	// if object exist, but no expiration is set on it then send -1
	if obj.ExpiresAt == -1 {
		c.Write([]byte(":-1\r\n"))
		return nil
	}

	// compute the time remaining for the key to expire and
	// return the RESP encoded form of it
	durationMs := obj.ExpiresAt - time.Now().UnixMilli()

	// if key expired i.e. key does not exist hence return -2
	if durationMs < 0 {
		c.Write([]byte(":-2\r\n"))
		return nil
	}

	c.Write(Encode(int64(durationMs/1000), false))
	return nil
}
func evalDEL(args []string, c io.ReadWriter) error {
	var countDeleted int = 0
	for _, key := range args {
		if ok := Del(key); ok {
			countDeleted++
		}
	}
	c.Write(Encode(countDeleted, false))
	return nil
}
func evalExpire(args []string, c io.ReadWriter) error {

	if len(args) <= 1 {
		return errors.New("(error) ERR wrong number of arguments for 'expire' command")
	}

	var key string = args[0]
	exDurationSec, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		return errors.New("(error) ERR value is not an integer or out of range")
	}

	obj := Get(key)

	// 0 if the timeout was not set. e.g. key doesn't exist, or operation skipped due to the provided arguments
	if obj == nil {
		c.Write([]byte(":0\r\n"))
		return nil
	}

	obj.ExpiresAt = time.Now().UnixMilli() + exDurationSec*1000

	// 1 if the timeout was set.
	c.Write([]byte(":1\r\n"))
	return nil
}

func EvaluateandRespond(cmd *Rediscmd, c io.ReadWriter) error {
	log.Println(cmd.Cmd)

	switch cmd.Cmd {
	case "PING":
		return evalPing(cmd.Args, c)
	case "SET":
		return evalSet(cmd.Args, c)
	case "GET":
		return evalGet(cmd.Args, c)
	case "TTL":
		return evalTTL(cmd.Args, c)
	case "DEL":
		return evalDEL(cmd.Args, c)
	case "EXPIRE":
		return evalExpire(cmd.Args, c)
	default:
		return evalPing(cmd.Args, c)

	}

}
