// Package plexer provides functionalities to manage a multiplexed connection.
package plexer

/*
// Multi represents a multiplexed connection to a server.
type Multi struct {
	connected bool   // True if a connection has been established, false otherwise
	service   string // The service to connect to
	address   string // The address of the server
	protocol  string // The network protocol to use (e.g., TCP, UDP)
	id        string // A unique identifier for this connection

	request net.Conn // The underlying network connection
}

// NewMulti initializes a new Multi instance and connects it to a server.
// address: Server's address
// service: Service to be connected to
// protocol: Network protocol to be used
// Returns a pointer to the new Multi instance, or an error if one occurs.
func NewMulti(address, service, protocol string) (*Multi, error) {
	v4, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	mp := &Multi{
		address:  address,
		service:  service,
		protocol: protocol,
		id:       v4.String(),
	}

	if err := mp.Connect(); err != nil {
		return nil, err
	}

	return mp, nil
}

// Connect establishes a connection to the server.
// Returns an error if the connection fails or if already connected.
func (mp *Multi) Connect() error {
	var err error

	if mp.connected {
		return fmt.Errorf("already connected")
	}

	mp.request, err = net.DialTimeout(mp.protocol, mp.address+":"+mp.service, time.Second*config.DEFAULT_TIMEOUT)
	if err != nil {
		return fmt.Errorf("can't establish connection")
	}

	mp.connected = true

	go mp.handleServerResponses()

	return err
}

// Disconnect closes the connection.
// Returns an error if the disconnection fails.
func (mp *Multi) Disconnect() error {
	mp.connected = false

	if err := mp.request.Close(); err != nil {
		mp.connected = true
		return err
	}

	return nil
}

// Write sends data over the connection.
// data: The data to send
// Returns the number of bytes written and any error encountered.
func (mp *Multi) Write(data []byte) (int, error) {
	return mp.request.Write(data)
}

// handleServerResponses listens for responses from the server and processes them.
func (mp *Multi) handleServerResponses() {
	reader := bufio.NewReader(mp.request)
	for {
		if !mp.connected {
			fmt.Println("break:!mp.connected")
			break
		}

		var length uint32
		if err := binary.Read(reader, binary.LittleEndian, &length); err != nil {
			fmt.Println("break:binary.Read")
			break
		}

		data := make([]byte, length)
		_, err := io.ReadFull(reader, data)
		if err != nil {
			mp.Disconnect()
			if err == io.EOF {
				logger.Info("connection closed by the server.")
			} else {
				logger.Error(err)
			}
			break
		}

		res := &transport.Response{}
		err = proto.Unmarshal(data, res)
		if logger.Error(err) {
			logger.Info(data)
			continue
		}

		if res.Pid != "" {
			p, err := promise.Find(res.Pid)
			if err != nil {
				fmt.Println(err)
				continue
			}

			p.Resolve(res)
		}
	}
}
*/
