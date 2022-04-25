package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path"
	"strings"
)

type Request struct {
	method string
	path   string
}

type Response struct {
	code        int
	contentType string
	body        string
}

func main() {
	listen, err := net.Listen("tcp", ":8080")
	checkError(err)
	fmt.Println("Server started.")
	for {
		conn, err := listen.Accept()
		checkError(err)
		fmt.Printf("Accepted connection from: %s\n", conn.RemoteAddr())
		inputBuffer := make([]byte, 1024)
		_, err = conn.Read(inputBuffer)
		checkError(err)
		request := parseRequest(inputBuffer)
		path := convertRelativePath(request.path)
		fmt.Printf("Served %s to %s\n", path, conn.RemoteAddr())
		serveFile(conn, path)
	}
}

func parseRequest(inputBuffer []byte) Request {
	lines := strings.Split(string(inputBuffer), "\n")
	return Request{
		method: strings.Fields(lines[0])[0],
		path:   strings.Fields(lines[0])[1],
	}
}

func convertRelativePath(relative string) string {
	workingDir, err := os.Getwd()
	checkError(err)
	webDir := path.Join(workingDir, "www")
	fullPath := path.Join(webDir, path.Clean(relative))
	return fullPath
}

func serveFile(conn net.Conn, filePath string) {
	info, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		sendError(conn, 404)
		return
	}
	if info.IsDir() {
		filePath = path.Join(filePath, "index.html")
	}
	data, err := ioutil.ReadFile(filePath)
	if os.IsNotExist(err) {
		sendError(conn, 404)
	}
	checkError(err)
	response := Response{
		code:        200,
		contentType: "text/html",
		body:        string(data),
	}
	sendHTML(conn, response)
}

func sendError(conn net.Conn, errorCode int) {
	sendHTML(conn, Response{
		code: errorCode,
		body: fmt.Sprintf("<b>Error: %v</b><br /><marquee>You fucked up!</marquee>", errorCode),
	})
}

func sendHTML(conn net.Conn, response Response) {
	header := fmt.Sprintf("HTTP/1.1 %v OK\r\n"+
		"Server: some go program\r\n"+
		"Connection: close\r\n"+
		"Content-Type: %s; charset=UTF-8\r\n"+
		"Content-Length: %v\r\n\r\n", response.code, response.contentType, len(response.body))
	output := header + response.body
	conn.Write([]byte(output))
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
