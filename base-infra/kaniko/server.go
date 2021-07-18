package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const (
	Message       = "Image built."
	StopCharacter = "\r\n\r\n"
)

func SocketServer(port int) {

	listen, err := net.Listen("tcp4", ":"+strconv.Itoa(port))

	if err != nil {
		log.Fatalf("Socket listen port %d failed,%s", port, err)
		os.Exit(1)
	}

	defer listen.Close()

	log.Printf("Begin listen port: %d", port)

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Fatalln(err)
			continue
		}
		go handler(conn)
	}

}

func handler(conn net.Conn) {

	defer conn.Close()

	var (
		buf  = make([]byte, 1024)
		r    = bufio.NewReader(conn)
		w    = bufio.NewWriter(conn)
		data = ""
	)

ILOOP:
	for {
		n, err := r.Read(buf)
		data = string(buf[:n])

		switch err {
		case io.EOF:
			break ILOOP
		case nil:
			log.Println("Receive:", data)
			if isTransportOver(data) {
				break ILOOP
			}
		default:
			log.Fatalf("Receive data failed:%s", err)
			return
		}
	}
	data = strings.TrimSuffix(data, "\r\n\r\n")
	log.Println("DATA:", data)
	args := strings.Split(data, "|")
	log.Println("ARGS:", args)
	projectName := args[0]
	containerRegistry := args[1]
	environment := args[2]
	app := "/executor"
	arg1 := fmt.Sprintf("--dockerfile=/workspace/%s/Dockerfile", projectName)
	arg2 := fmt.Sprintf("--context=dir://workspace/%s", projectName)
	arg3 := fmt.Sprintf("--destination=%s/%s:%s", containerRegistry, projectName, environment)
	cmd := exec.Command(app, arg1, arg2, arg3)
	cmd.Dir = "/"
	stdout, err := cmd.Output()

	if err != nil {
		log.Println("Error building image.")
		log.Println(err.Error())
		return
	}
	log.Println("Image built successfully.")
	// Print the output
	log.Println(string(stdout))

	w.Write([]byte(Message))
	w.Flush()
	log.Printf("Send: %s", Message)
}

func isTransportOver(data string) (over bool) {
	over = strings.HasSuffix(data, "\r\n\r\n")
	return
}

func main() {

	port := 3333

	SocketServer(port)

}
