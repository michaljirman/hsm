package main

import (
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

const (
	message       = "Ping"
	StopCharacter = "\r\n\r\n"
)

func SocketClient(ip string, port int) {
	addr := strings.Join([]string{ip, strconv.Itoa(port)}, ":")
	conn, err := net.Dial("tcp", addr)

	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}

	defer conn.Close()

	conn.Write([]byte(message))
	conn.Write([]byte(StopCharacter))
	log.Printf("Send: %s", message)

	buff := make([]byte, 1024)
	n, _ := conn.Read(buff)
	log.Printf("Receive: %s", buff[:n])
}

func main() {
	var (
		ip   = "localhost"
		port = 8098
	)
	SocketClient(ip, port)
}
