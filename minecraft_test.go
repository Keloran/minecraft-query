package minecraft_test

import (
	"fmt"
	"testing"

	minecraft "github.com/keloran/minecraft-query"
	"github.com/stretchr/testify/assert"
)

func TestConnectTCP(t *testing.T) {
	tests := []struct {
		name string
		mc   minecraft.Minecraft
		err  error
	}{
		{
			name: "LocalHost",
			mc: minecraft.Minecraft{
				Address: "localhost",
				Port:    25565,
				Timeout: 10,
			},
		},
		{
			name: "Failed",
			mc: minecraft.Minecraft{
				Address: "localhost",
				Port:    9999,
				Timeout: 10,
			},
			err: fmt.Errorf("invalid address"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mc, err := test.mc.ConnectTCP()
			passed := assert.IsType(t, test.err, err)
			if !passed {
				t.Errorf("test err: %w, request: %v", err, test)
			}

			if mc.Conn != nil {
				mc.Disconnect()
			}
		})
	}
}

func TestConnectUDP(t *testing.T) {
	tests := []struct {
		name string
		mc   minecraft.Minecraft
		err  error
	}{
		{
			name: "LocalHost",
			mc: minecraft.Minecraft{
				Address: "localhost",
				Port:    25565,
				Timeout: 10,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mc, err := test.mc.ConnectUDP()
			passed := assert.IsType(t, test.err, err)
			if !passed {
				t.Errorf("test err: %w, request: %v", err, test)
			}

			if mc.Conn != nil {
				mc.Disconnect()
			}
		})
	}
}
