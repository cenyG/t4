package model

import "fmt"

type StorageServerAddress struct {
	Port int    `json:"port"`
	Host string `json:"host"`
}

func (s StorageServerAddress) String() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}
