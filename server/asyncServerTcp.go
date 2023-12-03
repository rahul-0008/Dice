package server

import (
	"log"
	"net"
	"syscall"
	"time"

	"github.com/DiceDB/Dice/config"
	"github.com/DiceDB/Dice/core"
)

var con_clients = 0
var cronFrequency time.Duration = 1 * time.Second
var lastExecuted time.Time = time.Now()

func RunAsyncTCPServer() error {
	log.Println("Starting a Asynchronous Server on", config.Host, config.Port)

	// max clients that can be monitored
	max_clients := 20000

	// create a Kqueue event aobject array to store kQueue events
	var kqueue_events []syscall.Kevent_t = make([]syscall.Kevent_t, max_clients)

	// create a socket
	serverFD, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		return err
	}
	defer syscall.Close(serverFD)

	// set that sockt non-blockingg
	if err = syscall.SetNonblock(serverFD, true); err != nil {
		return err
	}

	//Bind the ip
	ipv4 := net.ParseIP(config.Host)
	if err = syscall.Bind(serverFD, &syscall.SockaddrInet4{
		Port: config.Port,
		Addr: [4]byte{ipv4[0], ipv4[1], ipv4[2], ipv4[3]},
	}); err != nil {
		return err
	}

	// started listening
	if err = syscall.Listen(serverFD, max_clients); err != nil {
		return err
	}
	//AsyncIO starts here!

	kq, err := syscall.Kqueue()
	if err != nil {
		return err
	}

	defer syscall.Close(kq)

	serverEvent := syscall.Kevent_t{
		Ident:  uint64(serverFD),
		Flags:  syscall.EV_ADD,
		Filter: syscall.EVFILT_READ,
	}
	//Listen to events on server FD
	_, err = syscall.Kevent(kq, []syscall.Kevent_t{serverEvent}, nil, nil)
	if err != nil {
		return err
	}

	for {

		if time.Now().After(lastExecuted.Add(cronFrequency)) {
			core.DeleteExpiredKeys()
			lastExecuted = time.Now()
		}
		// see and wait for events to happen in the kernel buffer and user space
		// this is a blocking call and to timeout with a specific time fill in 4th argument.

		nEvents, err := syscall.Kevent(kq, nil, kqueue_events[:], nil)
		if err != nil {
			continue
		}

		for i := 0; i < nEvents; i++ {
			event := kqueue_events[i]
			tempFD := int(event.Ident)

			// sevrver socket is ready for IO
			if tempFD == serverFD {
				clientFD, _, err := syscall.Accept(tempFD)
				if err != nil {
					log.Println("err : ", err)
					continue
				}
				con_clients++
				syscall.SetNonblock(serverFD, true)

				// add the new client TCP to be monitored
				tempClienEvent := syscall.Kevent_t{
					Ident:  uint64(clientFD),
					Flags:  syscall.EV_ADD,
					Filter: syscall.EVFILT_READ,
				}

				_, err = syscall.Kevent(kq, []syscall.Kevent_t{tempClienEvent}, nil, nil)
				if err != nil {
					log.Fatal(err)
				}

			} else {
				// this is one of the connected clients

				comm := core.FDComm{Fd: tempFD}
				cmds, err := readCommands(comm)
				if err != nil {
					syscall.Close(comm.Fd)
					con_clients -= 1
					// log.Println("Total clients connected now : ", con_clients)
					continue
				}
				respond(cmds, comm)

			}

		}

	}

}
