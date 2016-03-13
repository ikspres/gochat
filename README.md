# gochat
Simple chatting server and client written in Go


## Sample usage
(See /sample directory)


### Server
```
package main

import "github.com/ikspres/gochat/server"

func main() {
	cr := server.NewChatRoom(":6666")
	cr.Go()
```

### Client

```
package main

import (
	"github.com/ikspres/gochat/client"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <client nickname>", os.Args[0])
	}

	cli := client.NewClient(":6666", os.Args[1])
	defer cli.Close()

	cli.Go()
}

### running
```
go build cli.go
go build svr.go

# run server
svr

# run first client giving  nickname as argument
cli superman

# run second client
cli batman
```
