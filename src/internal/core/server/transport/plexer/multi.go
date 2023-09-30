package plexer

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/google/uuid"
	"github.com/kodmain/kitsune/src/config"
	"github.com/kodmain/kitsune/src/internal/core/server/promise"
	"github.com/kodmain/kitsune/src/internal/core/server/transport"
	"github.com/kodmain/kitsune/src/internal/kernel/observability/logger"
	"google.golang.org/protobuf/proto"
)

type Multi struct {
	connected bool

	service  string
	address  string
	protocol string
	id       string

	request net.Conn
	//response net.Conn
}

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

func (mp *Multi) Disconnect() error {
	mp.connected = false

	if err := mp.request.Close(); err != nil {
		mp.connected = true
		return err
	}

	return nil
}

func (mp *Multi) Write(data []byte) (int, error) {
	return mp.request.Write(data)
}

// handleServerResponses continuously reads responses from the server and forwards them to the appropriate channels.
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
