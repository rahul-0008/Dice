package server

import (
	"io"
	"log"
	"net"
	"strconv"

	"github.com/DiceDB/Dice/config"
)

func readCommand(c net.Conn) (string, error) {
	// TODO: Max read in one shot is 512 bytes
	// To allow input > 512 bytes, then repeated read until
	// we get EOF or designated delimiter

	var byt []byte = make([]byte, 512)
	n, err := c.Read(byt[:])

	if err != nil {
		return "", err
	}

	return string(byt[:n]), nil
}

func respond(cmd string, c net.Conn) error {
	if _, err := c.Write([]byte(cmd)); err != nil {
		return err
	}
	return nil
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

			log.Println("Command :", cmd)
			writeError := respond(cmd, c)
			if writeError != nil {
				log.Println("Write error", writeError)
			}
		}
	}
}
