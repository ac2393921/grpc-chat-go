package chat

import "net"

type room struct {
	name    string
	members map[net.Addr]*client
}
