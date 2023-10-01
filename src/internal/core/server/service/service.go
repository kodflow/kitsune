package service

import (
	"bytes"

	"github.com/kodmain/kitsune/src/internal/core/server/transport/plexer"
)

type Service struct {
	Name string
	mp   *plexer.Multi
}

func (s *Service) Write(data bytes.Buffer) (int, error) {
	return s.mp.Write(data.Bytes())
}

func (s *Service) Disconnect() error {
	return s.mp.Disconnect()
}

func (s *Service) MakeRequestWithResponse() *Query {
	return query(s.Name, true)
}

func (s *Service) MakeRequestOnly() *Query {
	return query(s.Name, false)
}

func Create(address, port, protocol string) (*Service, error) {
	mp, err := plexer.NewMulti(address, port, protocol)
	if err != nil {
		return nil, err
	}

	return &Service{
		Name: protocol + "//" + address + ":" + port,
		mp:   mp,
	}, nil
}
