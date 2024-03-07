package udpcommunication

import (
	"fmt"
	"net"
)

// UDPServer representa un servidor UDP
type UDPServer struct {
	conn *net.UDPConn
	outputChannel chan string
	inputChannel chan string
	sender *UDPClient
}

// NewUDPServer crea un nuevo servidor UDP
func NewUDPServer(port int, outputChannel chan string, inputChannel chan string, send *UDPClient) (*UDPServer, error) {
	udpAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return nil, err
	}

	return &UDPServer{
		conn: conn,
		outputChannel: outputChannel,
		inputChannel: inputChannel,
		sender: send,
	}, nil
}

func (s *UDPServer) SendMessage() {
	for message := range s.inputChannel{
		s.sender.SendMessage(message)
	}
}
// ReceiveMessage espera y recibe un mensaje del cliente UDP
func (s *UDPServer) ReceiveMessage() {
	buffer := make([]byte, 1024*100)
	fmt.Println("Longitud del buffer ", len(buffer))
	for true{
		
		n, _, err := s.conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Error reading message:", err)
			return
		}
		
		message := string(buffer[:n])
		if message != "ACK"{
			s.outputChannel <- message
			s.inputChannel <- "ACK"
		}
	}
}
