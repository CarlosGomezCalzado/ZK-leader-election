package udpcommunication

import (
	"fmt"
	"net"
)

// UDPClient representa un cliente UDP
type UDPClientAux struct {
	ServerAddr *net.UDPAddr
	conn       *net.UDPConn
}

type UDPClient struct{
	Nodes []int
}


func NewUDPC ()(*UDPClient, error){
	emptySlice := []int{}
	return &UDPClient{
		Nodes:emptySlice,
	},nil
}
// NewUDPClient crea un nuevo cliente UDP
func NewUDPClient(serverAddr string) (*UDPClientAux, error) {
	udpAddr, err := net.ResolveUDPAddr("udp", serverAddr)
	if err != nil {
		return nil, err
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		return nil, err
	}

	return &UDPClientAux{
		ServerAddr: udpAddr,
		conn:       conn,
	}, nil
}

// SendMessage envía un mensaje al servidor UDP
// NewUDPClient crea un nuevo cliente UDP
func (uc *UDPClient) AddNode (id int){
	uc.Nodes = append(uc.Nodes,id)
}

// SendMessage envía un mensaje al servidor UDP
func (uc *UDPClient) SendMessage(message string) {
	for _, id := range uc.Nodes {
		addr := fmt.Sprintf("127.0.0.1:%d",8080+id)
		client1, err := NewUDPClient(addr)
		if err != nil {
			fmt.Println("Error creating UDP client:", err)
			return
		}
		client1.SendMessage(message)
	}
	
}
func (c *UDPClientAux) SendMessage(message string) {
	data := []byte(message)
	_, err := c.conn.Write(data)
	if err != nil {
		fmt.Println("Error sending message:", err)
	}
	c.conn.Close()
}