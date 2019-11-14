package minecraft

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
)

func (m Minecraft) Ping(resolve bool) (ServerInfo, error) {
	if m.Conn == nil {
		return ServerInfo{}, fmt.Errorf("no connection")
	}

	if resolve {
		err := m.Resolve()
		if err != nil {
			return ServerInfo{}, fmt.Errorf("Ping Resolve: %w", err)
		}
	}

	return m.pingQuery()
}

func (m Minecraft) pingQuery() (ServerInfo, error) {
	var buf, buff bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, []byte("\x00\x04"))
	binary.Write(&buf, binary.LittleEndian, []byte(string(len(m.Address))))
	binary.Write(&buf, binary.LittleEndian, []byte(m.Address))
	binary.Write(&buf, binary.BigEndian, uint16(m.Port))
	binary.Write(&buf, binary.LittleEndian, []byte("\x01"))
	binary.Write(&buff, binary.LittleEndian, []byte(string(buf.Len())))
	binary.Write(&buff, binary.LittleEndian, buf.Bytes())
	binary.Write(&buff, binary.LittleEndian, []byte("\x01\x00"))

	_, err := m.Conn.Write(buff.Bytes())
	if err != nil {
		return ServerInfo{}, fmt.Errorf("handshake write: %w", err)
	}
	buff.Reset()
	buf.Reset()

	r := bufio.NewReader(m.Conn)
	z, err := r.ReadBytes('\n')
	if err != nil {
		if err != io.EOF {
			return ServerInfo{}, fmt.Errorf("read byte buffer: %w", err)
		}
	}

	j := ServerInfo{}
	err = json.Unmarshal(z[3:], &j)
	if err != nil {
		return ServerInfo{}, fmt.Errorf("unmarshal json: %w", err)
	}

	return j, nil
}

