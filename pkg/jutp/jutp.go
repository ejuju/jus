package jutp

import (
	"bufio"
	"io"
	"log"
	"net"
)

type Message string

func Write(w io.Writer, msg Message) (int, error) {
	return w.Write(append([]byte(msg), 0))
}

func Read(r *bufio.Reader) (Message, error) {
	msg, err := r.ReadBytes(0)
	if err != nil {
		return "", err
	}
	return Message(msg[:len(msg)-1]), nil
}

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
		go handler(&RemoteUI{conn: conn, r: bufio.NewReader(conn)})
	}
}

type RemoteUI struct {
	conn *net.TCPConn
	r    *bufio.Reader
}

func (rui *RemoteUI) Exec(code string) error { _, err := Write(rui.conn, Message(code)); return err }
func (rui *RemoteUI) Read() (Message, error) { return Read(rui.r) }
