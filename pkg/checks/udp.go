package checks

import (
	"encoding/json"
	"fmt"
	"github.com/aabri-ankorstore/cli-auth/pkg/server/util/port"
	"github.com/aabri-ankorstore/cli-auth/pkg/utils"
	"github.com/pkg/errors"
	"net"
	"os"
	"path/filepath"
	"strconv"
)

const host = "localhost"

type Payload struct {
	IsAuthenticated bool `json:"is_authenticated"`
}

type UdpProtocol struct {
	Host         string
	ForcePort    *int
	DefaultPort  int
	Conn         *net.UDPConn
	PluginFolder string
}

func NewUdpProtocol(host string, defaultPort int, pluginFolder string) *UdpProtocol {
	return &UdpProtocol{
		Host:         host,
		DefaultPort:  defaultPort,
		PluginFolder: pluginFolder,
	}
}

func (u *UdpProtocol) Listen() {
	usePort := u.DefaultPort
	if u.ForcePort != nil {
		usePort = *u.ForcePort
		if u.Host == host {
			available, err := port.IsAvailable(fmt.Sprintf(":%d", usePort))
			if !available {
				u.CheckError(errors.Errorf("Port %d already in use: %v", usePort, err))
			}
		}
	} else {
		if u.Host == host {
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
	u.CheckError(err)
	u.Conn, err = net.ListenUDP("udp", udpAddr)
	u.CheckError(err)
	fmt.Println("Udp Server started...")
	for {
		// Receive data
		u.HandleClient()
	}
}

func (u *UdpProtocol) IsAuth() bool {
	file := fmt.Sprintf("%s/%s/%s", u.PluginFolder, utils.PluginPath, pattern)
	matches, err := filepath.Glob(file)
	u.CheckError(err)
	if len(matches) > 0 {
		return true
	}
	return false
}

func (u *UdpProtocol) CheckError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error ", err.Error())
		os.Exit(1)
	}
}

func (u *UdpProtocol) HandleClient() {
	var buf [512]byte
	n, addr, err := u.Conn.ReadFromUDP(buf[0:])
	u.CheckError(err)
	var p Payload
	err = json.Unmarshal(buf[0:n], &p)
	u.CheckError(err)

	fmt.Println(p.IsAuthenticated)
	var jsonStr = []byte(fmt.Sprintf(`{"is_authenticated":"%t"}`, u.IsAuth()))
	_, err = u.Conn.WriteToUDP(jsonStr, addr)
	u.CheckError(err)
}

func Client() bool {
	udpAddr, err := net.ResolveUDPAddr("udp4", ":1200")
	if err != nil {
		panic(err)
	}
	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		panic(err)
	}

	// Send Data
	go func() {
		var jsonStr = []byte(`{"is_authenticated":"true"}`)
		_, err := conn.Write(jsonStr)
		if err != nil {
			panic(err)
		}
	}()

	for {
		// Receive Data
		go func() bool {
			var buf [512]byte
			n, err := conn.Read(buf[0:])
			if err != nil {
				panic(err)
			}

			var p Payload
			err = json.Unmarshal(buf[0:n], &p)
			if err != nil {
				panic(err)
			}
			return p.IsAuthenticated
		}()
	}
}
