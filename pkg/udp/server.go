package udp

import (
	"fmt"
	"github.com/aabri-ankorstore/cli-auth/pkg/filesystem"
	"github.com/aabri-ankorstore/cli-auth/pkg/port"
	"github.com/aabri-ankorstore/cli-auth/pkg/utils"
	"github.com/pkg/errors"
	"net"
	"os"
	"path/filepath"
	"strconv"
)

const host = "localhost"

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
	for {
		// Receive data
		u.HandleClient()
	}
}

func (u *UdpProtocol) IsAuth() bool {
	file := fmt.Sprintf("%s/%s/%s", u.PluginFolder, utils.PluginPath, filesystem.Pattern)
	matches, err := filepath.Glob(file)
	u.CheckError(err)
	if len(matches) > 0 {
		return true
	}
	return false
}

func (u *UdpProtocol) HandleClient() {
	var buf [512]byte
	_, addr, err := u.Conn.ReadFromUDP(buf[0:])
	u.CheckError(err)
	_, err = u.Conn.WriteToUDP([]byte(fmt.Sprintf("%t", u.IsAuth())), addr)
	u.CheckError(err)
}

func (u *UdpProtocol) CheckError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error ", err.Error())
		os.Exit(1)
	}
}
