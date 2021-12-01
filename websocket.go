package easyRpc

// WebsocketConn .
type WebsocketConn interface {
	HandleWebsocket(func())
}
