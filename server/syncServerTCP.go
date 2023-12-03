package server

import (
	"io"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/DiceDB/Dice/config"
	"github.com/DiceDB/Dice/core"
)

func convertToArrayString(ai []interface{}) ([]string, error) {
	as := make([]string, len(ai))
	for i := range ai {
		as[i] = ai[i].(string)
	}
	return as, nil
}

func readCommands(c io.ReadWriter) (core.RedisCmds, error) {
	// TODO: Max read in one shot is 512 bytes
	// To allow input > 512 bytes, then repeated read until
	// we get EOF or designated delimiter

	var byt []byte = make([]byte, 512)
	n, err := c.Read(byt[:])

	if err != nil {
		return nil, err
	}
	values, err := core.Decode(byt[:n])
	if err != nil {
		return nil, err
	}

	var cmds []*core.Rediscmd = make([]*core.Rediscmd, 0)
	for _, value := range values {
		tokens, err := convertToArrayString(value.([]interface{}))
		if err != nil {
			return nil, err
		}

		cmds = append(cmds, &core.Rediscmd{
			Cmd:  strings.ToUpper(tokens[0]),
			Args: tokens[1:],
		})
	}
	return cmds, nil
}

func respond(cmds core.RedisCmds, c io.ReadWriter) {
	core.EvaluateandRespond(cmds, c)
}

func RunSyncTCPServer() {
	log.Println("Starting a Synchronous Server on", config.Host, config.Port)

	var con_clients = 0

	lsnr, err := net.Listen("tcp", config.Host+":"+strconv.Itoa(config.Port))
	if err != nil {
		panic(err)
	}

	for {
		// this is a blocking call: it waits here for a new CLient to connect

		c, err := lsnr.Accept()
		if err != nil {
			panic(err)
		}

		con_clients += 1
		log.Println("client connected with address : ", c.RemoteAddr(), "concurrent clients", con_clients)

		for {
			// over the socket, continously read the command and print out

			cmds, err := readCommands(c)

			if err != nil {
				c.Close()
				con_clients -= 1
				log.Println("clients Disconnected", c.RemoteAddr(), "concurrent clients", con_clients)
				if err == io.EOF {
					break
				}
				log.Println("err", err)
			}
			respond(cmds, c)
		}
	}
}
