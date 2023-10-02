package jutp

import "net"

type Client struct {
	stream *net.TCPConn
}

func NewClient(laddr, raddr *net.TCPAddr) (*Client, error) {
	conn, err := net.DialTCP("tcp", laddr, raddr)
	if err != nil {
		return nil, err
	}
	return &Client{stream: conn}, nil
}

func (c *Client) Close() error { return c.stream.Close() }

func (c *Client) Send(msg string) error {
	_, err := c.stream.Write(append([]byte(msg), 0))
	if err != nil {
		return err
	}
	return nil
}
