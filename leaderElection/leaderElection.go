package leaderElection

import (
	"fmt"
	//a "ZK-leader-election/lthash"
	a "ZK-leader-election/udpcommunication"
	"encoding/json"
	"encoding/hex"
	"strconv"
	"crypto/elliptic"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"time"
)

// LeaderElection is the interface defining the "test" method.
type LeaderElection interface {
	StartServer(*a.UDPClient)
	Received() bool
	From() string
	Send(message string)
	Leader() int
}

// LElection represents the implementation of the LeaderElection interface.
type LElection struct{
	id int
	server *a.UDPServer
	from string
	received bool
	messageReceptionChannel chan string
	messageSendingChannel chan string
	born []string
	leader int
	bornHash string
	publicKey *ecdsa.PublicKey
	privateKey *ecdsa.PrivateKey
}

// MessageStruct represents the structure of the JSON messages exchanged.
type MessageStruct struct {
	Type 	string 	`json:"type"`
	ID 		int 	`json:"id"`
	Length	int 	`json:"length"`
	Hash	string 	`json:"hash"`
}

// hexStringToBoolSlice converts a hexadecimal string to a boolean slice.
func hexStringToBoolSlice(hexString string) ([]bool, error) {
	// Decode the hexadecimal string to []byte
	byteSlice, err := hex.DecodeString(hexString)
	if err != nil {
		return nil, err
	}

	// Convert []byte to []bool
	var boolSlice []bool
	for _, b := range byteSlice {
		boolSlice = append(boolSlice, b != 0x00)
	}

	return boolSlice, nil
}

// boolSliceToHexString converts a boolean slice to a hexadecimal string.
func boolSliceToHexString(boolSlice []bool) (string, error) {
	// Convert []bool to []byte
	var byteSlice []byte
	for _, b := range boolSlice {
		if b {
			byteSlice = append(byteSlice, 0x01)
		} else {
			byteSlice = append(byteSlice, 0x00)
		}
	}

	// Convert []byte to hexadecimal string
	hexString := hex.EncodeToString(byteSlice)
	return hexString, nil
}

// NewLE creates a new LeaderElection instance.
func NewLE(id int) LeaderElection{
	result := LElection{id:id, leader:id}
	return &result
}

// Leader returns the current leader.
func (le *LElection) Leader() int{
	return le.leader
}

// InBornList checks if a given string is in the born list.
func (le *LElection) InBornList(searchString string) bool {
	for _, s := range le.born {
		if s == searchString {
			return true
		}
	}
	return false
}

// handleBorn handles 'born' messages.
func (le *LElection) handleBorn(message MessageStruct) {
	if !le.InBornList(message.Hash) {
		le.born = append(le.born, message.Hash)
	} 
}

// handleLeader handles 'leader' messages.
func (le *LElection) handleLeader(message MessageStruct) {
	if len(le.born) == message.Length && le.leader < message.ID  {
		fmt.Println("I'm ", le.id, " and I change my leader to ", message.ID)
		le.leader = message.ID
	} else if len(le.born) < message.Length {
		fmt.Println("I'm ", le.id, " and I change my leader to ", message.ID)
		le.leader = message.ID
	}
}

// MessageReception listens for incoming messages and processes them based on their type.
func (le *LElection) MessageReception(){
	for message := range le.messageReceptionChannel {
		var messageObject MessageStruct
		err := json.Unmarshal([]byte(message), &messageObject)
		if err != nil {
			fmt.Println("Error decoding JSON:", err)
		}

		switch messageObject.Type {
		case "born":
			le.handleBorn(messageObject)
		case "leader":
			le.handleLeader(messageObject)
		default:
			fmt.Println("Unknown type:", messageObject.Type)
		}
	}
}

// LeaderRequest sends leader messages periodically.
func (le *LElection) LeaderRequest(){
	for true {
		if le.leader == le.id {
			message := MessageStruct{
				Type:    "leader",
				ID: le.id,
				Length: len(le.born),
			}
		
			// Encode the structure into a JSON string
			jsonMessage, err := json.Marshal(message)
			if err != nil {
				fmt.Println("Error encoding JSON:", err)
			}
			le.Send(string(jsonMessage))
		} else {
			fmt.Println("I'm ",le.id," and my leader is ",le.leader)
		}
		
		time.Sleep(200 * time.Millisecond)
	}
}

// Received returns whether a message has been received.
func (le *LElection) Received() bool{
	return le.received
}

// From returns the source of the message.
func (le *LElection) From() string{
	return le.from
}

// Send sends a message through the message sending channel.
func (le *LElection) Send(message string) {
	le.messageSendingChannel <- message
}

// StartServer initializes the leader election server.
func (le *LElection) StartServer(sender *a.UDPClient) {

	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		fmt.Println("Error generating private key:", err)
		return
	}
	le.privateKey = privateKey
	le.publicKey = &privateKey.PublicKey
	idString := strconv.Itoa(le.id)
	hashed := sha256.Sum256([]byte(idString))
	r, s, err := ecdsa.Sign(rand.Reader, le.privateKey, hashed[:])
	if err != nil {
		fmt.Println("Error signing the message:", err)
		return
	}
	rHex := fmt.Sprintf("%x", r)
	sHex := fmt.Sprintf("%x", s)

	// Combine r and s into a single string
	signatureHex := "0x" + rHex + sHex
	le.bornHash = signatureHex

	le.messageReceptionChannel = make(chan string)
	le.messageSendingChannel = make(chan string)
	le.born = append(le.born, signatureHex)
	server, err := a.NewUDPServer(8080+le.id, le.messageReceptionChannel, le.messageSendingChannel, sender)
	if err != nil {
		fmt.Println("Error creating UDP server:", err)
		return
	}
	go server.ReceiveMessage()
	go server.SendMessage()
	go le.MessageReception()
	go le.LeaderRequest()

	message := MessageStruct{
		Type:    "born",
		ID: le.id,
		Hash: signatureHex,
	}

	// Encode the structure into a JSON string
	jsonMessage, err := json.Marshal(message)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
	}
	le.Send(string(jsonMessage))

}
