package socket

import "net"

type client struct {
	conn *net.UnixConn
}

/*

func Client(socket string) (*client, error) {
	if _, err := os.Stat(socket); !os.IsNotExist(err) {
		return nil, err
	}

	conn, err := net.DialUnix(ADDR_TYPE, nil, addr(socket))

	if err != nil {
		return nil, err
	}

	return &client{conn}, nil
}

func (c *client) Send(Request *query.Request) (int, error) {
	bytes, err := json.Marshal(query.NewMessage(Request))

	if err != nil {
		return 0, err
	}

	return c.conn.Write(bytes)
}
*/
