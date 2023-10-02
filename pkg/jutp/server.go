package jutp

import (
	"log"
	"net"
)

func Serve(laddr *net.TCPAddr, handler func(rui *RemoteUI)) error {
	l, err := net.ListenTCP("tcp", laddr)
	if err != nil {
		return err
	}
	defer l.Close()

	for {
		conn, err := l.AcceptTCP()
		if err != nil {
			log.Println(err)
			continue
		}
		go handler(&RemoteUI{conn: conn})
	}
}

type RemoteUI struct {
	conn *net.TCPConn
}

func (rui *RemoteUI) Exec(code string) error {
	panic("todo")
}

func (rui *RemoteUI) Read() (string, error) {
	panic("todo")
}
