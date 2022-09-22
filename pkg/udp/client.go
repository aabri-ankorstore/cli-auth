package udp

import (
	"fmt"
	"net"
	"strconv"
)

func Client() (bool, error) {
	errChan := make(chan error, 0)
	results := make(chan string, 0)
	udpAddr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf("%s:%d", "localhost", 1200))
	if err != nil {
		errChan <- err
	}
	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		errChan <- err
	}
	// Send Data
	go func() {
		_, err := conn.Write([]byte(`ping`))
		if err != nil {
			errChan <- err
		}
	}()
	for {
		// Receive Data
		go func() {
			var buf [512]byte
			n, err := conn.Read(buf[0:])
			if err != nil {
				errChan <- err
			}
			results <- string(buf[0:n])
		}()
		select {
		case err := <-errChan:
			r, _ := strconv.ParseBool("false")
			return r, err
		case res := <-results:
			r, _ := strconv.ParseBool(res)
			return r, nil
			//os.Exit(0)
		}
	}
}
