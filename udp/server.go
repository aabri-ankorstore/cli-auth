package udp

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
)

type Payload struct {
	Message string `json:"message"`
}

func Run() {
	udpAddr, err := net.ResolveUDPAddr("udp4", ":1200")
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
