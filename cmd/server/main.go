package main

import (
	"bufio"
	"log"
	"net"
	"url-shortener/internal/handler"
	"url-shortener/internal/router"
	"url-shortener/internal/storage"
	"url-shortener/internal/web"
)

func main() {
	store := storage.NewMemoryStorage()
	h := handler.New(store)
	r := router.New(h)

	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}
		go handleConnection(conn, r)
	}

}

func handleConnection(conn net.Conn, r *router.Router) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	req, err := web.ReadRequest(reader)
	if err != nil {
		log.Printf("Failed to read request: %v", err)
		return
	}
	response := r.Handle(req)
	data := web.WriteResponse(response)
	_, err = conn.Write(data)
	if err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}
