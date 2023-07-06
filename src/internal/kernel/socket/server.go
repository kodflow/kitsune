package socket

import "net"

type server struct {
	socketPath string
	listener   *net.UnixListener
}

func (s *server) Start() error {
	return nil
}

func (s *server) Stop() error {
	return nil
}

/*
func Server(socketPath string) (*server, error) {
	os.Remove(socketPath)

	listener, err := net.ListenUnix("unixpacket", &net.UnixAddr{
		Name: socketPath,
		Net:  "unixpacket",
	})

	if err != nil {
		return nil, fmt.Errorf("erreur lors de la création du serveur: %w", err)
	}

	server := &server{
		socketPath: socketPath,
		listener:   listener,
	}

	return server, nil
}

func (s *server) Accept() ([]byte, error) {
	conn, err := s.listener.AcceptUnix()
	if err != nil {
		return nil, fmt.Errorf("erreur lors de l'acceptation de la connexion: %w", err)
	}
	defer conn.Close()

	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la réception des données du client: %w", err)
	}

	message := buffer[:n]
	return message, nil
}

func (s *server) Respond(message []byte) error {
	conn, err := s.listener.AcceptUnix()
	if err != nil {
		return fmt.Errorf("erreur lors de l'acceptation de la connexion: %w", err)
	}
	defer conn.Close()

	_, err = conn.Write(message)
	if err != nil {
		return fmt.Errorf("erreur lors de l'envoi de la réponse au client: %w", err)
	}

	return nil
}

func (s *server) Close() error {
	err := s.listener.Close()
	if err != nil {
		return fmt.Errorf("erreur lors de la fermeture du serveur: %w", err)
	}

	os.Remove(s.socketPath)

	return nil
}

*/
