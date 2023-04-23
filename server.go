package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"strings"
	"time"
)

var random *rand.Rand

func main() {
	s1 := rand.NewSource(time.Now().UnixNano())
	random = rand.New(s1)
	startServer()
}

func startServer() {
	udpServer, err := net.ListenPacket("udp", ":1054")
	if err != nil {
		log.Fatal(err)
	}
	defer udpServer.Close()
	requestHandler(udpServer)
}

func requestHandler(server net.PacketConn) {
	for {
		buf := make([]byte, 1024)
		_, addr, err := server.ReadFrom(buf)
		if err != nil {
			continue
		}
		if !buffLenValidator(buf, 100) {
			fmt.Print("Size of Buff is more than limitation!")
			return
		}
		go responseHandler(server, addr, buf)
	}
}

func responseHandler(udpServer net.PacketConn, addr net.Addr, buf []byte) {
	if responseProb := random.Intn(10); responseProb > 5 {
		return
	}
	response := responseCreator(string(buf))
	udpServer.WriteTo([]byte(response), addr)
	fmt.Printf("sent %s to client\n", response)
}

func responseCreator(input string) string {
	return strings.ReplaceAll(input, "bali", "kheyr")
}

func buffLenValidator(buff []byte, limit int) bool {
	length := 0
	for _, byt := range buff {
		if length > limit {
			return false
		}
		if byt == 0 {
			break
		}
	}
	return true
}
