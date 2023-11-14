package server

import (
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/DiceDB/Dice/config"
	"github.com/DiceDB/Dice/core"
)

func readCommand(c io.ReadWriter) (*core.Rediscmd, error) {
	// TODO: Max read in one shot is 512 bytes
	// To allow input > 512 bytes, then repeated read until
	// we get EOF or designated delimiter

	var byt []byte = make([]byte, 512)
	n, err := c.Read(byt[:])

	if err != nil {
		return nil, err
	}
	tokens, err := core.DecodeArrayString(byt[:n])
	if err != nil {
		return nil, err
	}

	return &core.Rediscmd{
		Cmd:  strings.ToUpper(tokens[0]),
		Args: tokens[1:],
	}, nil
}

func respondError(err error, c io.ReadWriter) {

	b := []byte(fmt.Sprintf("-%s\r\n", err))
	c.Write(b)
}

func respond(cmd *core.Rediscmd, c io.ReadWriter) {
	err := core.EvaluateandRespond(cmd, c)

	if err != nil {
		respondError(err, c)
	}
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

			cmd, err := readCommand(c)

			if err != nil {
				c.Close()
				con_clients -= 1
				log.Println("clients Disconnected", c.RemoteAddr(), "concurrent clients", con_clients)
				if err == io.EOF {
					break
				}
				log.Println("err", err)
			}
			respond(cmd, c)
		}
	}
}
