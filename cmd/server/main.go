package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"url-shortener/internal/handler"
	"url-shortener/internal/router"
	"url-shortener/internal/storage"
	"url-shortener/internal/web"
)

func main() {
	store, err := newStorage()
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}
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

func newStorage() (storage.Storage, error) {
	switch driver := os.Getenv("STORAGE_DRIVER"); driver {
	case "", "memory":
		return storage.NewMemoryStorage(), nil
	case "postgres":
		dsn := os.Getenv("DATABASE_URL")
		if dsn == "" {
			return nil, fmt.Errorf("DATABASE_URL is required when STORAGE_DRIVER=postgres")
		}
		return storage.NewPostgresStorage(dsn)
	default:
		return nil, fmt.Errorf("unknown STORAGE_DRIVER: %q", driver)
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
