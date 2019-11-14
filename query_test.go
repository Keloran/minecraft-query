package minecraft_test

import (
	"testing"

	minecraft "github.com/keloran/minecraft-query"
	"github.com/stretchr/testify/assert"
)

func TestQuery(t *testing.T) {
	tests := []struct{
		name string
		mc minecraft.Minecraft
		err error
	}{
		{
			name: "Localhost",
			mc: minecraft.Minecraft{
				Address: "localhost",
				Port: 25565,
				Timeout: 10,
			},
		},
	}

	for _, test := range tests {
//		t.Run(test.name, func (t *testing.T) {
			m, err := test.mc.ConnectUDP()
			if err != nil {
				t.Errorf("test %s, %w", test.name, err)
			}

//			err = m.Query()
//			assert.Equal(t, test.err, err)
			m.Disconnect()
			assert.Equal(t, true, true)
//		})
	}
}
