package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

var (
	projectName       = os.Getenv("GOCD_PROJECT_NAME")
	containerRegistry = os.Getenv("GOCD_CONTAINER_REGISTRY")
	environment       = os.Getenv("GOCD_ENVIRONMENT")
	args              = fmt.Sprintf("%s|%s|%s", projectName, containerRegistry, environment)
	StopCharacter     = "\r\n\r\n"
)

func SocketClient(ip string, port int) {
	addr := strings.Join([]string{ip, strconv.Itoa(port)}, ":")
	conn, err := net.Dial("tcp", addr)

	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}

	defer conn.Close()

	conn.Write([]byte(args))
	conn.Write([]byte(StopCharacter))
	log.Printf("Send: %s", args)

	buff := make([]byte, 1024)
	n, _ := conn.Read(buff)
	log.Printf("Receive: %s", buff[:n])

}

func main() {

	var (
		ip   = "127.0.0.1"
		port = 3333
	)

	SocketClient(ip, port)

}
