package main

import (
	"./graph"
	"./message"
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"time"
)

const myConst = 4096

type task struct {
	message     message.Message
	destination net.UDPAddr
}

func runGossipNerwork(curID int, graph graph.Graph, quitChan, killChan chan struct{}, ttl int, needInit bool) {
	curNode, _ := graph.GetNode(curID)
	curPort := curNode.Port()
	curAddress, error := net.ResolveUDPAddr("udp", "127.0.0.1:"+strconv.Itoa(curPort))
	if error != nil {
		panic(error)
	}
	curConnection, error := net.ListenUDP("udp", curAddress)
	if error != nil {
		panic(error)
	}

	neighbours, _ := graph.Neighbors(curID)
	neighboursAmount := len(neighbours)

	neighboursPorts := make([]int, neighboursAmount)
	neighboursIDs := make([]int, neighboursAmount)
	neighboursAdresses := make([]net.UDPAddr, neighboursAmount)

	for i := range neighbours {
		neighboursPorts[i] = neighbours[i].Port()
		neighboursIDs[i], _ = strconv.Atoi(neighbours[i].String())

		address, error := net.ResolveUDPAddr("udp", "127.0.0.1:"+strconv.Itoa(neighboursPorts[i]))
		if error != nil {
			panic(error)
		}
		neighboursAdresses[i] = *address
	}

	dataChannel := make(chan message.Message)
	taskChannel := make(chan task, myConst)
	timeoutChannel := make(chan time.Duration, 1)
	timeoutChannel <- time.Millisecond * 50
	quitChannel1 := make(chan struct{}, 1)
	quitChannel2 := make(chan struct{}, 1)
	go manageConnection(curConnection, dataChannel, timeoutChannel, quitChannel1)
	go sendQueueMessages(curConnection, taskChannel, time.Microsecond*10, quitChannel2)
	time.Sleep(time.Millisecond * 10)

	nodesReceptionStatus := make([]bool, len(graph))
	if needInit {
		for i := 0; i < ttl; i++ {
			destination := rand.Intn(neighboursAmount)
			address := neighboursAdresses[destination]
			m := task{message.Message{ID: curID, Type: "message", Sender: 0, Origin: 0, Data: "data"}, address}
			taskChannel <- m
		}
		for index := range nodesReceptionStatus {
			nodesReceptionStatus[index] = false
		}
		nodesReceptionStatus[curID] = true
	}

	//ttlMap := make(map[int]int)
	ticks := 0
MessageProcessing:
	for {
		select {
		case message := <-dataChannel:
			{
				select {
				case <-killChan:
					{
					}
				default:
					{
						if message.Type == "message" {
							//MessageConfirm
						} else if message.Type == "confirmative" {
							if true {
								fmt.Println("Finished in", ticks, "ticks")
								break MessageProcessing
							}
						} else if message.Type != "timeout" {
							panic("Unknown message type" + message.Type)
						}
					}
				}
			}
		}
		ticks++
	}

	quitChannel1 <- struct{}{}
	quitChannel2 <- struct{}{}
	quitChan <- struct{}{}
}

func sendQueueMessages(curConnection *net.UDPConn, taskChannel chan task, interval time.Duration, quitChan chan struct{}) {
	var buf = make([]byte, myConst)

	for {
		select {
		case _ = <-quitChan:
			return
		case task := <-taskChannel:
			{
				time.Sleep(interval)
				buf = task.message.ConvertToJsonMsg()
				_, error := curConnection.WriteToUDP(buf, &task.destination)
				if error != nil {
					panic(error)
				}
			}
		}
	}
}

func manageConnection(Conn *net.UDPConn, dataChannel chan message.Message, timeoutChannel chan time.Duration, quitChan chan struct{}) {
	var buf = make([]byte, myConst)
	var msg message.Message

	timeout := time.Duration(0)
	for {
		select {
		case timeout = <-timeoutChannel:
			{
			}
		case _ = <-quitChan:
			{
				return
			}
		default:
			{
				if timeout != 0 {
					Conn.SetReadDeadline(time.Now().Add(timeout))
				}
				messageLength, _, error := Conn.ReadFromUDP(buf)

				if error != nil {
					err, ok := error.(net.Error)
					if !ok || !err.Timeout() {
						panic(error)
					}

					msg = message.Message{ID: 0, Type: "timeout", Sender: 0, Origin: 0, Data: ""}
				} else {
					msg = message.ConvertFromJsonMsg(buf[0:messageLength])
				}

				dataChannel <- msg
			}
		}
	}
}
