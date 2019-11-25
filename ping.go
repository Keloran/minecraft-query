package minecraft

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
)

func (m Minecraft) Ping(resolve bool) (PingInfo, error) {
	if m.Conn == nil {
		return PingInfo{}, fmt.Errorf("no connection")
	}

	if resolve {
		err := m.Resolve()
		if err != nil {
			return PingInfo{}, fmt.Errorf("Ping Resolve: %w", err)
		}
	}

	return m.pingQuery()
}

func (m Minecraft) pingQuery() (PingInfo, error) {
	var buf, buff bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, []byte{00, 04})
	if err != nil {
	    return PingInfo{}, fmt.Errorf("pingQuery magic: %w", err)
    }
	err = binary.Write(&buf, binary.LittleEndian, []byte(string(len(m.Address))))
	if err != nil {
	    return PingInfo{}, fmt.Errorf("pingQuery address length: %w", err)
    }
	err = binary.Write(&buf, binary.LittleEndian, []byte(m.Address))
	if err != nil {
	    return PingInfo{}, fmt.Errorf("pingQuery address: %w", err)
    }
	err = binary.Write(&buf, binary.BigEndian, uint16(m.Port))
	if err != nil {
	    return PingInfo{}, fmt.Errorf("pingQuery port: %w", err)
    }
	binary.Write(&buf, binary.LittleEndian, []byte{01})
	binary.Write(&buff, binary.LittleEndian, []byte(string(buf.Len())))
	binary.Write(&buff, binary.LittleEndian, buf.Bytes())
	binary.Write(&buff, binary.LittleEndian, []byte{01, 00})

	_, err = m.Conn.Write(buff.Bytes())
	if err != nil {
		return PingInfo{}, fmt.Errorf("handshake write: %w", err)
	}
	buff.Reset()
	buf.Reset()

	r := bufio.NewReader(m.Conn)
	z, err := r.ReadBytes('\n')
	if err != nil {
		if err != io.EOF {
			return PingInfo{}, fmt.Errorf("read byte buffer: %w", err)
		}
	}

	j := PingInfo{}
	err = json.Unmarshal(z[3:], &j)
	if err != nil {
		return PingInfo{}, fmt.Errorf("unmarshal json: %w", err)
	}

	return j, nil
}
