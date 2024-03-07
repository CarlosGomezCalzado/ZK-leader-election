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
	b "ZK-leader-election/lthash"
	"bloom"
)

// LeaderElection es la interfaz que define el método prueba
type LeaderElection interface {
	StartServer(*a.UDPClient)
	Received() bool
	From() string
	Send(message string)
	Leader() int
}

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
	LTHash 	   b.Hash
	Bloom 	*c.BloomFilter
}

type MessageStruct struct {
	Type 	string 	`json:"type"`
	ID 		int 	`json:"id"`
	Length	int 	`json:"length"`
	Hash	string 	`json:"hash"`
	LTHash 	[]byte 	`josn:"lthash"`
	Bloom	string 	`json:"bloom"`
}

func hexStringToBoolSlice(hexString string) ([]bool, error) {
	// Decodifica la cadena hexadecimal a []byte
	byteSlice, err := hex.DecodeString(hexString)
	if err != nil {
		return nil, err
	}

	// Convierte []byte a []bool
	var boolSlice []bool
	for _, b := range byteSlice {
		boolSlice = append(boolSlice, b != 0x00)
	}

	return boolSlice, nil
}

func boolSliceToHexString(boolSlice []bool) (string, error) {
	// Convierte []bool a []byte
	var byteSlice []byte
	for _, b := range boolSlice {
		if b {
			byteSlice = append(byteSlice, 0x01)
		} else {
			byteSlice = append(byteSlice, 0x00)
		}
	}

	// Convierte []byte a cadena hexadecimal
	hexString := hex.EncodeToString(byteSlice)
	return hexString, nil
}

func NewLE(id int) LeaderElection{
	result := LElection{id:id, leader:id}
	return &result
}

func (le *LElection) Leader() int{
	return le.leader
}

func (le *LElection) InBornList(searchString string) bool {
	for _, s := range le.born {
		if s == searchString {
			return true
		}
	}
	return false
}

func (le *LElection) handleBorn(message MessageStruct) {
	// fmt.Println(le.id, " Ejecutando función para el tipo 'born'. Contenido:", message.ID, message.Length)
	// Puedes agregar aquí la lógica específica para el tipo 'born'
	if !le.InBornList(message.Hash) {
		le.born = append(le.born, message.Hash)
		le.LTHash.Add([]byte(message.Hash))
		le.Bloom.Add([]byte(message.Hash))
	} 
}

func (le *LElection) handleLeader(message MessageStruct) {
	// fmt.Println(le.id, "Ejecutando función para el tipo 'leader'. Contenido:", message.ID, message.Length)
	h := b.New16()
	h.SetState(message.LTHash)
	if message.ID == le.id {
		return
	}
	boolSlice, err := hexStringToBoolSlice(message.Bloom)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	//fmt.Println("0")
	if h.In(le.born) == 0 && le.leader < message.ID  {
		fmt.Println("Soy ", le.id, " y cambio a ", message.ID)
		le.leader = message.ID
	} else if h.In(le.born) < 0 {
		// fmt.Println(h.In(le.born))
		bf := c.SetStateBloom(boolSlice)
		// fmt.Println(le.id)
		// fmt.Println(message.Bloom)
		if bf.Contains([]byte(le.bornHash)) {
			fmt.Println("Soy ", le.id, " y cambio a ", message.ID)
			le.leader = message.ID
		}
	}
	// fmt.Println("4")
	// if len(le.born) < message.Length {
	// 	le.leader = message.ID
	// } else if len(le.born) == message.Length && le.id < message.ID {
	// 	le.leader = message.ID
	// }
	// Puedes agregar aquí la lógica específica para el tipo 'leader'
}

func (le *LElection) MessageReception(){
	for message := range le.messageReceptionChannel {
		var messageObject MessageStruct
		err := json.Unmarshal([]byte(message), &messageObject)
		if err != nil {
			fmt.Println("Error al decodificar JSON:", err)
		}

		switch messageObject.Type {
		case "born":
			le.handleBorn(messageObject)
		case "leader":
			le.handleLeader(messageObject)
		default:
			fmt.Println("Tipo desconocido:", messageObject.Type, " es lo que llega de ", message)
		}
	}
}

func (le *LElection) LeaderRequest(){
	for true {
		if le.leader == le.id {
			hexString, err := boolSliceToHexString(le.Bloom.GetBitmap())
			if err != nil {
				fmt.Println("Error:", err)
			}
			message := MessageStruct{
				Type:    "leader",
				ID: le.id,
				Length: len(le.born),
				Bloom: hexString,
				LTHash: le.LTHash.Sum(nil),
			}
		
			// Codificar la estructura en una cadena JSON
			jsonMessage, err := json.Marshal(message)
			if err != nil {
				fmt.Println("Error al codificar JSON:", err)
			}
			le.Send(string(jsonMessage))
		} else {
			fmt.Println("Soy ",le.id,"Mi lider es ",le.leader)
		}
		// if len(le.born) == 1 {
		// 	message := MessageStruct{
		// 		Type:    "born",
		// 		ID: le.id,
		// 		Hash: le.bornHash,
		// 	}
		
		// 	// Codificar la estructura en una cadena JSON
		// 	jsonMessage, err := json.Marshal(message)
		// 	if err != nil {
		// 		fmt.Println("Error al codificar JSON:", err)
		// 	}
		// 	le.Send(string(jsonMessage))
		// }
		time.Sleep(200 * time.Millisecond)
	}
}

func (le *LElection) Received() bool{
	return le.received
}

func (le *LElection) From() string{
	return le.from
}

func (le *LElection) Send(message string) {
	le.messageSendingChannel <- message
}

func (le *LElection) StartServer(sender *a.UDPClient) {

	le.LTHash = b.New16()
	le.Bloom = bloom.NewWithEstimates(1000000, 0.01) 
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

	// Combinar r y s en un solo string
	signatureHex := "0x" + rHex + sHex
	le.bornHash = signatureHex

	le.messageReceptionChannel = make(chan string)
	le.messageSendingChannel = make(chan string)
	le.born = append(le.born, signatureHex)
	le.LTHash.Add([]byte(signatureHex))
	le.Bloom.Add([]byte(signatureHex))
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

	// Codificar la estructura en una cadena JSON
	jsonMessage, err := json.Marshal(message)
	if err != nil {
		fmt.Println("Error al codificar JSON:", err)
	}
	le.Send(string(jsonMessage))

}
