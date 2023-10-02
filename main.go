package main

import (
	"net"

	"github.com/ejuju/jus/pkg/jul"
	"github.com/ejuju/jus/pkg/jutp"
)

func main() {
	// Connect to server
	client, err := jutp.NewClient(nil, &net.TCPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: 8080,
	})
	if err != nil {
		panic(err)
	}

	// Execute code received from server
	for {
		vm := jul.NewVM(jul.WithClient(client))
		err := vm.Execute(panic("todo"))
		if err != nil {
			panic(err)
		}
	}
}
