package easyRpc

// WebsocketConn defines websocket-conn interface.
type WebsocketConn interface {
	HandleWebsocket(func())
}
