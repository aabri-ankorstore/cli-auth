package udp

import (
	"encoding/json"
	"fmt"
	"github.com/aabri-ankorstore/cli-auth/server/util/port"
	"github.com/pkg/errors"
	"net"
	"os"
	"strconv"
)

type Payload struct {
	Message string `json:"message"`
}

// DefaultPort is the default port the ui server will listen to
const DefaultPort = 1200

func NewServer(host string, forcePort *int) {
	// Find an open port
	usePort := DefaultPort
	if forcePort != nil {
		usePort = *forcePort
		if host == "localhost" {
			available, err := port.IsAvailable(fmt.Sprintf(":%d", usePort))
			if !available {
				checkError(errors.Errorf("Port %d already in use: %v", usePort, err))
			}
		}
	} else {
		if host == "localhost" {
			for i := 0; i < 20; i++ {
				available, _ := port.IsAvailable(fmt.Sprintf(":%d", usePort))
				if available {
					break
				}
				usePort++
			}
		}
	}

	udpAddr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf("%s:%s", host, strconv.Itoa(usePort)))
	checkError(err)
	conn, err := net.ListenUDP("udp", udpAddr)
	checkError(err)
	for {
		handleClient(conn)
	}
}
func handleClient(conn *net.UDPConn) {
	var buf [512]byte
	n, addr, err := conn.ReadFromUDP(buf[0:])
	if err != nil {
		return
	}
	var p Payload
	err = json.Unmarshal(buf[0:n], &p)
	if err != nil {
		return
	}

	fmt.Println(p.Message)

	var jsonStr = []byte(`{"message":"pong"}`)
	conn.WriteToUDP(jsonStr, addr)
}
func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error ", err.Error())
		os.Exit(1)
	}
}
