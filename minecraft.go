package minecraft

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"time"
)

type Minecraft struct {
	Socket  string
	Address string
	Port    int
	Timeout int
	Conn    net.Conn
	Session int
}

type Description struct {
	Name string `json:"text"`
}
type Players struct {
	Max    int `json:"max"`
	Online int `json:"online"`
}
type Version struct {
	Version  string `json:"name"`
	Protocol int    `json:"protocol"`
}

type Plugin struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type Player struct {
	Name string `json:"name"`
}

type PingInfo struct {
	Description `json:"description"`
	Players     `json:"players"`
	Version     `json:"version"`
}

type HostInfo struct {
	Name string `json:"name"`
	Port int    `json:"port"`
	IP   string `json:"ip"`
}

type GameInfo struct {
	Type string `json:"type"`
	Name string `json:"name"`
	Map  string `json:"map"`
}

type PlayerInfo struct {
	Max    int `json:"max"`
	Online int `json:"online"`
}

type ServerInfo struct {
	HostInfo   `json:"host"`
	Version    `json:"version"`
	PlayerInfo `json:"player"`
}

type QueryInfo struct {
	ServerInfo `json:"server"`
	GameInfo   `json:"game"`
	Players    []Player `json:"players"`
	Plugins    []Plugin `json:"plugins"`
}

func (m Minecraft) Resolve() error {
	if ip2long(m.Address) == 0 {
		return fmt.Errorf("Invalid IP")
	}

	return nil
}

func ip2long(ip string) uint32 {
	var long uint32
	err := binary.Read(bytes.NewBuffer(net.ParseIP(ip).To4()), binary.BigEndian, &long)
	if err != nil {
		fmt.Printf("ip2long read: %v", err)
	}
	return long
}

func (m Minecraft) ConnectTCP() (Minecraft, error) {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", m.Address, m.Port), time.Duration(m.Timeout)*time.Second)
	if err != nil {
		return m, fmt.Errorf("connection timeout")
	}
	m.Conn = conn

	return m, nil
}

func (m Minecraft) ConnectUDP() (Minecraft, error) {
	conn, err := net.DialTimeout("udp", fmt.Sprintf("%s:%d", m.Address, m.Port), time.Duration(m.Timeout)*time.Second)
	if err != nil {
		fmt.Printf("connect: %v, %v\n", m, err)
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
