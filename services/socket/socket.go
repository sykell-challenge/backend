package socket

import (
	"fmt"
	"sync"

	"github.com/zishang520/engine.io/v2/types"
	"github.com/zishang520/socket.io/v2/socket"
)

var (
	globalServer *socket.Server
	serverMutex  sync.RWMutex
)

func InitSocketServer() *socket.Server {
	opts := socket.DefaultServerOptions()
	opts.SetCors(&types.Cors{
		Origin:      "*",
		Credentials: true,
	})
	opts.SetTransports(types.NewSet("polling", "websocket"))

	server := socket.NewServer(nil, opts)

	server.On("connection", func(clients ...interface{}) {
		client := clients[0].(*socket.Socket)
		fmt.Println("Client connected:", client.Id())

		// Join client to a general room for broadcasts
		client.Join("crawl_updates")
	})

	// Store the server globally for broadcasting
	serverMutex.Lock()
	globalServer = server
	serverMutex.Unlock()

	return server
}

// BroadcastCrawlUpdate broadcasts crawl job completion to all connected clients
func BroadcastCrawlUpdate(eventType string, data interface{}) {
	serverMutex.RLock()
	server := globalServer
	serverMutex.RUnlock()

	if server != nil {
		server.To("crawl_updates").Emit(eventType, data)
		fmt.Printf("Broadcasted %s event to all clients\n", eventType)
	} else {
		fmt.Println("Socket server not initialized, cannot broadcast")
	}
}

// GetServer returns the global socket server instance
func GetServer() *socket.Server {
	serverMutex.RLock()
	defer serverMutex.RUnlock()
	return globalServer
}
