package minecraft

import (
	"bytes"
	"encoding/binary"
	//	"encoding/hex"
	"fmt"
	"io"
	"strings"
	"strconv"
)

func (m Minecraft) Query() (interface{}, error) {
	if m.Conn == nil {
		return nil, fmt.Errorf("no connection")
	}

	challenge, err := m.GetChallenge()
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	status, err := m.GetStatus(challenge)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
//	fmt.Println(fmt.Sprintf("Status: %v", status))

	return status, nil
}

func (m Minecraft) GetChallenge() (int32, error) {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, []byte("\xFE\xFD"))
	binary.Write(&buf, binary.BigEndian, []byte("\x09"))
	binary.Write(&buf, binary.BigEndian, uint32(1))
	_, err := m.Conn.Write(buf.Bytes())
	if err != nil {
		return 0, fmt.Errorf("challange write: %w", err)
	}
	// fmt.Println(fmt.Sprintf(hex.EncodeToString(buf.Bytes())))

	out := make([]byte, 12)
	_, err = io.ReadFull(m.Conn, out)
	if err != nil {
		if err != io.EOF {
			return 0, fmt.Errorf("challenge read: %w", err)
		}
	}

	preChal := out[5:]
	chal := preChal
	if chal[len(chal) - 1] == byte(00) {
		chal = chal[0:len(chal) - 2]
	}

	ret, err := strconv.Atoi(string(chal))
	if err != nil {
		return 0, fmt.Errorf("challenge str to int: %w", err)
	}

	return int32(ret), nil
}

func (m Minecraft) GetStatus(challenge int32) (interface{}, error) {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, []byte("\xFE\xFD"))
	binary.Write(&buf, binary.BigEndian, []byte("\x00"))
	binary.Write(&buf, binary.BigEndian, uint32(1))
	binary.Write(&buf, binary.BigEndian, challenge)
	binary.Write(&buf, binary.BigEndian, []byte("\x00\x00\x00\x00"))
	_, err := m.Conn.Write(buf.Bytes())
	if err != nil {
		return "", fmt.Errorf("status write: %w", err)
	}
//	fmt.Println(fmt.Sprintf("status write: %x", buf.Bytes()))

	out := make([]byte, 1024)
	_, err = io.ReadAtLeast(m.Conn, out, 5)
	if err != nil {
		if err != io.EOF {
			return nil, fmt.Errorf("Status Read: %w", err)
		}
	}

	si, err := getServerInfo(out[16:])
	if err != nil {
		return nil, fmt.Errorf("status serverinfo: %w", err)
	}

//	fmt.Println(fmt.Sprintf("si: %+v", si))

	return si, nil
}

func getServerInfo(stat []byte) (QueryInfo, error) {
	ret := []byte{}
	data := []string{}

	for _, i := range stat {
		if i != 0 {
			ret = append(ret, i)
		} else {
			data = append(data, string(ret))
			ret = []byte{}
		}
	}

	q := QueryInfo{}

	for i, info := range data {
		switch info {
			// HostInfo
			case "hostname":
				q.ServerInfo.HostInfo.Name = data[i+1]
			case "hostport":
				port, _ := strconv.Atoi(data[i+1])
				q.ServerInfo.HostInfo.Port = port
			case "hostip":
				q.ServerInfo.HostInfo.IP = data[i+1]

			// Version
			case "version":
				q.ServerInfo.Version.Version = data[i+1]

			// Plugins
			case "plugins":
				plugins, err := getPlugins(stat)
				if err != nil {
					return q, fmt.Errorf("serverinfo plugins: %w", err)
				}
				q.Plugins = plugins

			// PlayerInfo
			case "numplayers":
				online, _ := strconv.Atoi(data[i+1])
				q.ServerInfo.PlayerInfo.Online = online
			case "maxplayers":
				max, _ := strconv.Atoi(data[i+1])
				q.ServerInfo.PlayerInfo.Max = max

			// GameInfo
			case "gametype":
				q.GameInfo.Type = data[i+1]
			case "game_id":
				q.GameInfo.Name = data[i+1]
			case "map":
				q.GameInfo.Map = data[i+1]

			// Players
			case "player_":
				if len(info) != 0 {
					players, err := getPlayers(stat)
					if err != nil {
						return q, fmt.Errorf("serverinfo players: %w", err)
					}
					q.Players = players
				}
		}

		if strings.Contains(info, "player_") {
			players, err := getPlayers(stat)
			if err != nil {
				return q, fmt.Errorf("serverinfo players: %w", err)
			}
			q.Players = players
		}
	}

//	fmt.Println(fmt.Sprintf("data: %v", data))
//	fmt.Println(fmt.Sprintf("Q: %+v", q))

	return q, nil
}

func getPlugins(stat []byte) ([]Plugin, error) {
	pr := []Plugin{}

	pluginBytes := []byte{112, 108, 117, 103, 105, 110, 115}
	postBytes := []byte{109, 97, 112}

	pluginFrom := 0
	pluginEnd := 0

	for n, i := range stat {
		if (i == pluginBytes[0] &&
		stat[n] == pluginBytes[0] &&
		stat[n+1] == pluginBytes[1] &&
		stat[n+2] == pluginBytes[2] &&
		stat[n+3] == pluginBytes[3] &&
		stat[n+4] == pluginBytes[4] &&
		stat[n+5] == pluginBytes[5] &&
		stat[n+6] == pluginBytes[6]) {
			pluginFrom = n+7
		}

		if i == postBytes[0] {
			if stat[n] == postBytes[0] && stat[n+1] == postBytes[1] && stat[n+2] == postBytes[2] {
				pluginEnd = n-1
			}
		}
	}

	pluginBytes = stat[pluginFrom:pluginEnd]
	endbucket := []byte{84, 58, 32}

	for n, i := range pluginBytes {
		if (i == endbucket[0] &&
		pluginBytes[n+1] == endbucket[1] &&
		pluginBytes[n+2] == endbucket[2]) {
			pluginFrom = n+3
		}
	}
	pluginBytes = pluginBytes[pluginFrom:]

	plStart := 0
	plEnd := 0

	semi := byte(59)
	for n, i := range pluginBytes {
		if i == semi {
			plEnd = n-1

			pr = append(pr, getPluginNameAndVersion(pluginBytes[plStart:plEnd]))

			plStart = n+2
		}

		if n == len(pluginBytes) - 1 {
			pr = append(pr, getPluginNameAndVersion(pluginBytes[plStart:]))
		}
	}

	return pr, nil
}

func getPluginNameAndVersion(pl []byte) Plugin {
	p := Plugin{}

	versionStart := 0

	numBytes := []byte{30, 31, 32, 33, 34, 35, 36, 37, 38, 39}

	for n, i := range pl {
		for _, ii := range numBytes {
			if i == ii {
				versionStart = n
				break
			}
		}
	}

	p.Name = string(pl[0:versionStart])
	p.Version = string(pl[versionStart+1:])

	return p
}

func getPlayers(stat []byte) ([]Player, error) {
	playerBytes := []byte{112, 108, 97, 121, 101, 114, 95} 
	plr := []Player{}
	
	playerStart := 0
	for n, i := range stat {
		if (i == playerBytes[0] &&
		stat[n+1] == playerBytes[1] &&
		stat[n+2] == playerBytes[2] &&
		stat[n+3] == playerBytes[3] &&
		stat[n+4] == playerBytes[4] &&
		stat[n+5] == playerBytes[5] &&
		stat[n+6] == playerBytes[6]) {
			playerStart = n+7
			break
		}
	}

	playersString := string(stat[playerStart:])
	players := strings.Split(playersString, " ")
	for _, i := range players {
		if len(i) >= 1 {
			player := Player{
				Name: i,
			}
			plr = append(plr, player)
		}
	}

	return plr, nil
}

