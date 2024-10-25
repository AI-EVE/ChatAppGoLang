package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		fmt.Println("failed to connect")
		return
	}
	defer conn.Close()

	
	go func ()  {
		messageReader := bufio.NewReader(os.Stdin)

		for {
			message, err := messageReader.ReadString('\n')
			if err != nil {
				fmt.Println("failed to read a message..")
				continue
			}

			_, err = conn.Write([]byte(message))
			if err != nil {
				fmt.Println("Failed to write to the server..")
			}
		}	
	}()
	
	readBuffer := make([]byte, 1000)
	for {
		n, err := conn.Read(readBuffer)
		if err != nil {
			fmt.Println("failed reading from the server", err)
		}

		fmt.Println(string(readBuffer[:n]))
	}

}