package main

import (
	"fmt"
	"math/rand"
	"net"
	"os"
	"sync"
	"time"
)

func main() {

	rand.Seed(time.Now().UnixNano())

	resp := make(chan []byte, 1)

	connection := connectToServer()

	listener(resp, connection)

	//close the connection
	defer connection.Close()

	sendMessage(connection, resp)

}

func listener(resp chan []byte, conn net.Conn) {
	go func() {
		received := make([]byte, 2048)
		_, err := conn.Read(received)
		if err != nil {
			println("Read data failed:", err.Error())
			os.Exit(1)
		}
		resp <- received
	}()
}
func connectToServer() *net.UDPConn {
	udpServer, err := net.ResolveUDPAddr("udp", ":1054")

	if err != nil {
		println("ResolveUDPAddr failed:", err.Error())
		os.Exit(1)
	}

	conn, err := net.DialUDP("udp", nil, udpServer)
	if err != nil {
		println("Listen failed:", err.Error())
		os.Exit(1)
	}

	return conn
}

func sendMessage(conn *net.UDPConn, resp chan []byte) {
	now := time.Now()
	message := createRandomMessage()
	for {
		fmt.Printf("sending string %s to server \n", message)
		_, err := conn.Write([]byte(message))
		if err != nil {
			println("Write data failed:", err.Error())
			os.Exit(1)
		}
		responseReader(now, resp)
		time.Sleep(2 * time.Second)
	}
}

func responseReader(start time.Time, resp chan []byte) {
	timeout := 4

	var wg sync.WaitGroup

	wg.Add(1)

	ticker := time.NewTicker(1 * time.Second)
	go func() {
		for {
			select {
			case response := <-resp:
				fmt.Printf("recieved string %s within %s time \n", string(response), time.Since(start).String())
				os.Exit(0)
			case <-ticker.C:
				timeout--
				if timeout == 0 {
					fmt.Println("No response from server. timeout deadline exceeded. retrying...")
					wg.Done()
					return
				}
			}
		}
	}()

	wg.Wait()
}
func createRandomMessage() string {
	size := rand.Intn(20)
	s := RandStringRunes(size)
	message := ""
	for _, c := range s {
		if random := rand.Intn(10); random < 6 {
			message += "bali"
		}
		message += string(c)
	}
	return message
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
