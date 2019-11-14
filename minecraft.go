package minecraft

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"time"
	"errors"
)

type Minecraft struct {
	Socket string
	Address string
	Port int
	Timeout int
	Conn net.Conn
}

type Description struct {
	Name string `json:"text"`
}
type Players struct {
	Max int `json:"max"`
	Online int `json:"online"`
}
type Version struct {
	Version string `json:"name"`
	Protocol int `json:"protocol"`
}

type ServerInfo struct {
	Description `json:"description"`
	Players `json:"players"`
	Version `json:"version"`
}

func (m Minecraft) Resolve() error {
	if ip2long(m.Address) == 0 {
		return fmt.Errorf("Invalid IP")
	}

	// record := net.LookupSVR("_minecraft._tcp", "tcp", m.ServerAddress)
	// fmt.Println(record)

	return nil
}

func ip2long(ip string) uint32 {
	var long uint32
	binary.Read(bytes.NewBuffer(net.ParseIP(ip).To4()), binary.BigEndian, &long)
	return long
}

func (m Minecraft) ConnectTCP() (Minecraft, error) {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", m.Address, m.Port), time.Duration(m.Timeout) * time.Second)
	if err != nil {
		fmt.Println(fmt.Errorf("connect: %v, %w", m, err))
		return m, errors.New("invalid address")
	}
	m.Conn = conn

	return m, nil
}

func (m Minecraft) ConnectUDP() (Minecraft, error) {
	conn, err := net.DialTimeout("udp", fmt.Sprintf("%s:%d", m.Address, m.Port), time.Duration(m.Timeout) * time.Second)
	if err != nil {
		fmt.Println(fmt.Errorf("connect: %v, %w", m, err))
		return m, errors.New("invalid address")
	}
	m.Conn = conn

	return m, nil
}

func (m Minecraft) Disconnect() error {
	if m.Conn != nil {
		err := m.Conn.Close()
		if err != nil {
			return fmt.Errorf("disconnect: %w", err)
		}
	}

	return nil
}

