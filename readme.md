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
        
    ping(m)
    query(m)
  }
        
  func ping(m mc.Minecraft) {        
    ml, err := m.ConnectTCP()
    if err != nil {
      fmt.Printf("ping tcp: %v\n", err)
      return
    }
        
    r, err := ml.Ping(false)
    if err != nil {
      fmt.Printf("ping err: %v\n", err)
      return
    }
        
    ml.Disconnect()
    fmt.Printf("ping: %+v\n", r)
  }
        
  func query(m mc.Minecraft) {
    ml, err := m.ConnectUDP()
    if err != nil {
      fmt.Printf("query udp: %v\n", err)
      return
    }
        
    r, err := ml.Query()
    if err != nil {
      fmt.Printf("query err: %v\n", err)
      return
    }
        
    ml.Disconnect()
    fmt.Printf("query: %+v\n", r)
  }
```

# Output
```go
    ping: {
      Description: {
        Name: "A Minecraft Server",
      },
      Players: {
        Max: 20,
        Online: 1,
      },
      Version: {
        Version: "CraftBukkit 1.14.4", 
        Protocol: 498,
      }
    }

    query: {
      ServerInfo: {
        HostInfo: {
          Name: "A Minecraft Server", 
          Port: 25565,
          IP: "127.0.0.1",
        },
        Version: {
          Version: "1.14.4",
          Protocol: 0,
        },
        PlayerInfo: {
          Max: 20,
          Online: 1,
        },
      },
      GameInfo: {
        Type: "SMP",
        Name: "MINECRAFT",
        Map: "world",
      },
      Players:[{
        Name: "Keloran",
      }],
      Plugins:[{
        Name: "WorldEdit",
        Version: "7.0.1,61bc01",
      }, {
        Name: "WorldGuard", 
        Version: "7.0.1-SNAPSHOT,556b638",
      }],
    }
    
```
