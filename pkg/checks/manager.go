package checks

import (
	"fmt"
)

type CheckManager interface {
	CreateTmpFile() (string, error)
	IsAuthenticatedOffline() bool
	Listen() error
	RemoveAuth()
	HandleClient() error
	CheckError(err error)
}

func GetType(t string) (CheckManager, error) {
	switch t {
	case "udp":
		return NewUdpProtocol(), nil
	case "filesystem":
		return NewFilesystem(), nil
	default:
		return nil, fmt.Errorf("wrong type passed")
	}
}
