package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	listen, err := net.Listen("tcp", ":8080")
	checkError(err)
	for {
		conn, err := listen.Accept()
		checkError(err)
		fmt.Println("Connected")
		sendHTML(conn)
	}
}

func sendHTML(conn net.Conn) {
	body := "unchi 💩"
	header := fmt.Sprintf("HTTP/1.1 200 OK\r\n"+
		"Server: some go program\r\n"+
		"Connection: close\r\n"+
		"Content-Type: text/html; charset=UTF-8\r\n"+
		"Content-Length: %v\r\n\r\n", len(body))
	response := header + body
	fmt.Println(response)

	conn.Write([]byte(response))
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
