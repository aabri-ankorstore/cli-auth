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
	Message string `json:"message"`
}

type UdpProtocol struct {
	Host         string
	ForcePort    *int
	DefaultPort  int
	Conn         *net.UDPConn
	PluginFolder string
	Payload
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
		u.CheckError(u.HandleClient())
	}
}
func (u *UdpProtocol) IsAuthenticated() bool {
	udpAddr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf("%s:%d", u.Host, u.DefaultPort))
	u.CheckError(err)
	conn, err := net.DialUDP("udp", nil, udpAddr)
	u.CheckError(err)

	go func() {
		var jsonStr = []byte(`{"message":"ping"}`)
		_, err := conn.Write(jsonStr)
		u.CheckError(err)
	}()
	for {
		go func() {
			var buf [512]byte
			n, err := conn.Read(buf[0:])
			u.CheckError(err)

			var p Payload
			err = json.Unmarshal(buf[0:n], &p)
			if err != nil {
				return
			}
			fmt.Println(p.Message)
			os.Exit(0)
		}()
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

func (u *UdpProtocol) HandleClient() error {
	var buf [512]byte
	n, addr, err := u.Conn.ReadFromUDP(buf[0:])
	if err != nil {
		return err
	}
	var p Payload
	err = json.Unmarshal(buf[0:n], &p)
	if err != nil {
		return err
	}

	//fmt.Println(p.Message)
	var jsonStr = []byte(fmt.Sprintf(`{"message":"%t"}`, true))
	u.Conn.WriteToUDP(jsonStr, addr)
	return nil
}
