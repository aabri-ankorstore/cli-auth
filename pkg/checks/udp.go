package checks

import (
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
	_, addr, err := u.Conn.ReadFromUDP(buf[0:])
	u.CheckError(err)
	_, err = u.Conn.WriteToUDP([]byte(fmt.Sprintf("%t", u.IsAuth())), addr)
	u.CheckError(err)
}

func Client() (string, error) {
	udpAddr, err := net.ResolveUDPAddr("udp4", ":1200")
	if err != nil {
		panic(err)
	}
	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		panic(err)
	}
	errChan := make(chan error, 0)
	results := make(chan string, 0)
	// Send Data
	go func() {
		_, err := conn.Write([]byte(`ping`))
		errChan <- err
	}()

	for {
		// Receive Data
		go func() {
			var buf [512]byte
			n, err := conn.Read(buf[0:])
			errChan <- err
			bytes := buf[0:n]
			results <- string(bytes)
			os.Exit(0)
		}()
		select {
		case err := <-errChan:
			return "false", err
		case res := <-results:
			return res, nil
		}
	}
}
