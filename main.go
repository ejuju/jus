package main

import (
	"bufio"
	"log"
	"net"
	"strings"

	"github.com/ejuju/jus/pkg/jul"
	"github.com/ejuju/jus/pkg/jutp"
)

func main() {
	// Connect to remote server
	conn, err := net.DialTCP("tcp", nil, &net.TCPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: 8080,
	})
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// Execute code received from server
	vm := jul.NewVM(jul.WithServerConnection(conn))
	r := bufio.NewReader(conn)
	msg, err := jutp.Read(r)
	if err != nil {
		log.Println(err)
		return
	}
	err = vm.Execute(strings.NewReader(string(msg)))
	if err != nil {
		log.Println(err)
		return
	}
}
