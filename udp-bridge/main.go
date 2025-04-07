package main

import (
	"fmt"
	"net"
	// "unsafe"
	// "golang.org/x/sys/windows"
)

const (
	shmName    = "Local\\acpmf_physics"  // Replace with your ACC SHM name
	SERVER_UDP_ADDR = "<localIp>:1234"
)

func ReadSharedMemory() {

}

func SharedMemory() {
	conn, err := net.Dial("udp", SERVER_UDP_ADDR)
	if err != nil {
		fmt.Println("[ERROR] - Connecting to UDP Server: ", err)
	}
	defer conn.Close()

	b := make([]byte, 1024)

	for {
		_, err = conn.Write(b)
		if err != nil {
			fmt.Println("[ERROR] - Writing to UDP Server: ", err)
		}
	}
}

func main() {
	fmt.Println("--- Running... ---")

	go func() {
		ReadSharedMemory()
	}()

	go func() {
		SharedMemory()
	}()
}