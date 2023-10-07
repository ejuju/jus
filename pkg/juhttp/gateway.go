package juhttp

import "net"

type Gateway struct {
	forwardTo *net.TCPAddr
}

func (gw *Gateway) Run(port int) {

}
