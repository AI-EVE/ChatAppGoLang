package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
)

var clients = make(map[net.Conn]string)
var mutex = sync.Mutex{}
func main() {
	listner, err := net.Listen("tcp", "127.0.0.1:8080")
    if err != nil {
        log.Println("couldnt listen to the server")
        return
    }
    defer listner.Close()

    closeChan := make(chan os.Signal, 1)
    signal.Notify(closeChan, os.Interrupt, syscall.SIGTERM)
    go func ()  {
        <- closeChan

        fmt.Println("Shutting down the server...")
        listner.Close()
        os.Exit(0)
    }()


    for {
        conn, err := listner.Accept()
        if err != nil {
            log.Println("A connection has failed")
            continue
        }

        go handleConnection(conn)
    }
}

func handleConnection(conn net.Conn) {
    defer conn.Close()
    mutex.Lock()
    userNumber := strconv.Itoa(len(clients) + 1)
    clients[conn] = userNumber
    for client := range clients {
        if client != conn {
            client.Write([]byte("User " + userNumber + "has entered the chat room...."))
        }
    }
    mutex.Unlock()

    buffer := make([]byte, 1024)
    for {
        n, err := conn.Read(buffer)
        if err != nil {
            fmt.Println("Failed to read...")
            break
        }

        if n == 0 {
            break
        }

        message := "User " + userNumber + ": " + string(buffer[:n])

        mutex.Lock()
        for client := range clients {
            if client != conn {
                client.Write([]byte(message))
            }
        }
        mutex.Unlock()
    }

    mutex.Lock()
    delete(clients, conn)
    for client := range clients {
        client.Write([]byte("User " + userNumber + " " + "has left.."))
    }
    mutex.Unlock()
} 