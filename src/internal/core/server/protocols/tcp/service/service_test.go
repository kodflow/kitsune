package service

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func startMockTCPServer(t *testing.T) func() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		t.Fatalf("Failed to start mock server: %v", err)
	}

	go func() {
		for {
			_, err := listener.Accept()
			if err != nil {
				return
			}
		}
	}()

	return func() {
		listener.Close()
	}
}

func TestCreate(t *testing.T) {
	stopServer := startMockTCPServer(t)
	defer stopServer()
	address := "127.0.0.1"
	port := "8080"

	svc, err := Create(address, port)
	assert.NoError(t, err)
	assert.NotNil(t, svc)
}

func TestConnectDisconnect(t *testing.T) {
	stopServer := startMockTCPServer(t)
	defer stopServer()
	address := "127.0.0.1"
	port := "8080"

	svc, _ := Create(address, port)
	assert.True(t, svc.Connected)
	svc.Connect()
	assert.True(t, svc.Connected)

	svc2, err := Create(address, "8081")
	assert.Nil(t, svc2)
	assert.Error(t, err)

	err = svc.Disconnect()
	assert.NoError(t, err)
	err = svc.Disconnect()
	assert.Error(t, err)
	assert.False(t, svc.Connected)
}

// MockConn simule l'interface net.Conn pour les tests.
type MockConn struct {
	Buffer    bytes.Buffer
	Connected bool
}

func (mc *MockConn) Read(b []byte) (n int, err error) {
	if !mc.Connected {
		return 0, errors.New("not connected")
	}
	return mc.Buffer.Read(b)
}

func (mc *MockConn) Write(b []byte) (n int, err error) {
	if !mc.Connected {
		return 0, errors.New("not connected")
	}
	return mc.Buffer.Write(b)
}

func (mc *MockConn) Close() error {
	mc.Connected = false
	return nil
}

func (mc *MockConn) LocalAddr() net.Addr                { return nil }
func (mc *MockConn) RemoteAddr() net.Addr               { return nil }
func (mc *MockConn) SetDeadline(t time.Time) error      { return nil }
func (mc *MockConn) SetReadDeadline(t time.Time) error  { return nil }
func (mc *MockConn) SetWriteDeadline(t time.Time) error { return nil }

func TestWrite(t *testing.T) {
	mockConn := &MockConn{Connected: true}
	svc := &Service{
		Name:      "localhost:8080",
		Address:   "localhost",
		Protocol:  "tcp",
		Network:   mockConn,
		Connected: true,
	}

	data := bytes.NewBufferString("test data")
	_, err := svc.Write(data)
	if err != nil {
		t.Fatalf("Write returned an error: %v", err)
	}

	if mockConn.Buffer.String() != "test data" {
		t.Errorf("Expected buffer to contain 'test data', got '%s'", mockConn.Buffer.String())
	}
}

func TestMakeExchange(t *testing.T) {
	svc := &Service{Name: "testService"}

	exchange1 := svc.MakeExchange()
	assert.NotNil(t, exchange1, "MakeExchange returned nil, expected non-nil Exchange")
	assert.Equal(t, svc.Name, exchange1.Service, "Service name does not match")
	assert.True(t, exchange1.Answer, "Expected Answer to be true, got false")

	exchange2 := svc.MakeExchange(false)
	assert.NotNil(t, exchange2, "MakeExchange with false returned nil, expected non-nil Exchange")
	assert.Equal(t, svc.Name, exchange2.Service, "Service name does not match")
	assert.False(t, exchange2.Answer, "Expected Answer to be false, got true")
}
func TestReadData(t *testing.T) {
	svc := &Service{Name: "testService"}

	reader := bufio.NewReader(bytes.NewBufferString("test data"))
	length := uint32(9)

	data, err := svc.readData(reader, length)
	assert.NoError(t, err)
	assert.Equal(t, []byte("test data"), data)

}

func TestHandleReadError(t *testing.T) {
	mockConn := &MockConn{Connected: true}
	svc := &Service{
		Name:      "localhost:8080",
		Address:   "localhost",
		Protocol:  "tcp",
		Network:   mockConn,
		Connected: true,
	}

	err := errors.New("test error")

	// Test case 1: Error is io.EOF
	err1 := io.EOF
	handleReadError(svc, err1)
	assert.False(t, svc.Connected, "Expected Connected to be false")

	// Test case 2: Error is not io.EOF
	handleReadError(svc, err)
	assert.False(t, svc.Connected, "Expected Connected to be false")
}
func TestUnmarshalResponse(t *testing.T) {
	data := []byte{ /* your test data here */ }

	res, err := unmarshalResponse(data)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	// Add your assertions for the unmarshaled response here
}
