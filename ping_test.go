package minecraft_test

import (
	"testing"

	minecraft "github.com/keloran/minecraft-query"
	"github.com/stretchr/testify/assert"
)

func TestPing(t *testing.T) {
	tests := []struct{
		name string
		mc minecraft.Minecraft
		expect minecraft.ServerInfo
		err error
	}{
		{
			name: "query",
			mc: minecraft.Minecraft{
				Address: "localhost",
				Port: 25565,
				Timeout: 10,
			},
			expect: minecraft.ServerInfo{
				Description: minecraft.Description{
					Name: "A Minecraft Server",
				},
				Players: minecraft.Players{
					Max: 20,
					Online: 0,
				},
				Version: minecraft.Version{
					Version: "1.14.4",
					Protocol: 498,
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			m, err := test.mc.ConnectTCP()
			if err != nil {
				t.Errorf("connect failed: %w", err)
			}

			ret, err := m.Ping(false)
			passed := assert.IsType(t, test.err, err)
			if !passed {
				t.Errorf("query: %w", err)
			}
			passed = assert.Equal(t, test.expect, ret)
			if !passed {
				test.mc.Disconnect()
				t.Errorf("test: %v", test)
			}
		})
	}
}

