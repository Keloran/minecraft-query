# Minecraft Query
This is a go port of the minecraft-query library from xPaw [PHP-Minecraft-Query](https://github.com/xPaw/PHP-Minecraft-Query)

# Example
```go
	package main
	import (
		"fmt"

		mc "github.com/keloran/minecraft-query"
	)

	func main() {
		m := mc.Minecraft{
			Address: "localhost",
			Port: 25565,
			Timeout: 10,
		}


		ml, err := m.ConnectUDP()
		if err != nil {
			fmt.Printf("Connect: %v\n", err)
			return
		}

		f := ml.Query()
		fmt.Printf("Query: %v\n", f)
		ml.Disconnect()
	}
```

# Output
```go
	QueryInfo: {
		ServerInfo: {
			HostInfo: {
				Name: "localhost",
				Port: 25565,
				IP: "127.0.0.1",
			},
			Version: "",
			PlayerInfo: {
				Max: 20,
				Online: 1,
			},
		},
		GameInfo: {
			Type: "",
			Name: "",
			Map: "",
		},
		Players: [{
			Name: "Keloran",
		}],
		Plugins: [{
			Name: "tester",
		}],
	}
```
