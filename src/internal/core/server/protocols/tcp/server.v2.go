// Package tcp provides functionalities for a TCP server.
package tcp

import (
	"bufio"
	"context"
	"encoding/binary"
	"errors"
	"io"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"syscall"

	"github.com/kodmain/kitsune/src/internal/core/server/api"
	"github.com/kodmain/kitsune/src/internal/core/server/router"
	"github.com/kodmain/kitsune/src/internal/core/server/transport"
	"github.com/kodmain/kitsune/src/internal/kernel/multithread"
	"github.com/kodmain/kitsune/src/internal/kernel/observability/logger"
	"golang.org/x/sys/unix"
	"google.golang.org/protobuf/proto"
)

//https://github.com/douglasmakey/socket-sharding/blob/master/cmd/http-example/main.go

// Server represents a TCP server and contains information about the address it listens on
// and the underlying network listener.
type ServerV2 struct {
	Address  string       // Address to listen on
	listener net.Listener // TCP Listener object
	config   net.ListenConfig
	childs   []*exec.Cmd
	running  bool
}

// NewServer creates a new Server instance with the specified listening address.
func NewServerV2(address string) *ServerV2 {
	return &ServerV2{
		Address: address,
		childs:  make([]*exec.Cmd, runtime.NumCPU()-1),
		config: net.ListenConfig{
			Control: func(network, address string, c syscall.RawConn) error {
				var opErr error
				if err := c.Control(func(fd uintptr) {
					opErr = unix.SetsockoptInt(int(fd), unix.SOL_SOCKET, unix.SO_REUSEPORT, 1)
				}); err != nil {
					return err
				}
				return opErr
			},
		},
	}
}

// Register is a method for registering API handlers with the server.
func (s *ServerV2) Register(api api.APInterface) {

}

// Start starts the TCP server, allowing it to accept incoming connections.
func (s *ServerV2) Start() error {
	if s.running {
		return errors.New("server already started")
	}

	s.running = true

	if multithread.IsMaster() {
		for i := range s.childs {
			s.childs[i] = exec.Command(os.Args[0], "-child")
			s.childs[i].Stdout = os.Stdout
			s.childs[i].Stderr = os.Stderr
			if err := s.childs[i].Start(); err != nil {
				return err
			}
		}

		for _, ch := range s.childs {
			if err := ch.Wait(); err != nil {
				return err
			}
		}

		logger.Info("master server start on " + s.Address + " with pid:" + strconv.Itoa(os.Getpid()))
	} else {
		runtime.GOMAXPROCS(1)
		var err error
		s.listener, err = s.config.Listen(context.Background(), "tcp", s.Address)
		if err != nil {
			return err
		}

		logger.Info("child server start on " + s.Address + " with pid:" + strconv.Itoa(os.Getpid()))
		go s.accepLoop()
	}

	return nil
}

// Stop stops the TCP server.
// Stop stops the server and kills all child processes.
// s: Server instance containing the necessary configurations and state.
// Returns: error if any operation fails.
func (s *ServerV2) Stop() error {
	if !s.running {
		return errors.New("server is not active")
	}

	s.running = false

	// If it's a master process, terminate all child processes
	if multithread.IsMaster() {
		for _, ch := range s.childs {
			if err := ch.Process.Kill(); err != nil {
				if err.Error() == "os: process already finished" {
					return nil
				}

				return err
			}
		}
		logger.Info("All child processes killed. Master server with pid: " + strconv.Itoa(os.Getpid()) + " is stopping.")
	} else {
		if err := s.listener.Close(); err != nil {
			return err
		}
		logger.Info("Child server on " + s.Address + " with pid: " + strconv.Itoa(os.Getpid()) + " is stopping.")
	}

	return nil
}

// acceptLoop continuously accepts incoming connections.
func (s *ServerV2) accepLoop() {
	for {
		if s.listener == nil {
			break
		}

		conn, err := s.listener.Accept()
		if err != nil {
			break
		}

		//logger.Info("new connection on " + s.Address + " from " + conn.RemoteAddr().Network() + " " + strconv.Itoa(os.Getpid()))

		go s.handleConnection(conn)
	}
}

// handleConnection handles incoming client connections.
func (s *ServerV2) handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	for {
		data, err := s.handleRequest(reader)
		if err != nil {
			break
		}

		s.sendResponse(writer, s.handler(data))
	}
}

// handleRequest reads and processes incoming requests.
func (s *ServerV2) handleRequest(reader *bufio.Reader) ([]byte, error) {
	var length uint32
	if err := binary.Read(reader, binary.LittleEndian, &length); err != nil {
		return nil, err
	}

	data := make([]byte, length)
	_, err := io.ReadFull(reader, data)
	return data, err
}

// sendResponse sends a response to the client.
func (s *ServerV2) sendResponse(writer *bufio.Writer, res []byte) {
	if len(res) > 0 {
		binary.Write(writer, binary.LittleEndian, uint32(len(res)))
		writer.Write(res)
		writer.Flush()
	}
}

// handler processes incoming data and generates a response.
func (s *ServerV2) handler(b []byte) []byte {
	req := &transport.Request{}
	res := &transport.Response{}
	err := proto.Unmarshal(b, req)
	if err != nil {
		return router.Empty
	}

	err = router.Resolve(req, res)
	if err != nil {
		return router.Empty
	}

	b, err = proto.Marshal(res)
	if err != nil {
		return router.Empty
	}

	return b
}
