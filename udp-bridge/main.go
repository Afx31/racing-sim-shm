package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"

	// "os"
	"sync"
	"time"
	"unsafe"

	"udp-bridge/internal/accshmdata"

	"github.com/hidez8891/shm"
)

const (
	shmNamePhysics    = "Local\\acpmf_physics"
	shmNameGraphics		= "Local\\acpmf_graphics"
	shmNameStatic			= "Local\\acpmf_static"
	SERVER_UDP_ADDR 	= "<localIp>:1234"
)

var (
	wg sync.WaitGroup
)

func ReadSharedMemory(physics *accshmdata.ACCPhysics) {
	fmt.Println("--- Reading Memory ---")
	
	// Setup ticker for looping beloiw
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		<-ticker.C
		physicsSize := (int32)(unsafe.Sizeof(*physics))
		rbuf := make([]byte, physicsSize)
		buf := &bytes.Buffer{}

		// Try open shared memory
		r, err := shm.Open(shmNamePhysics, physicsSize)
		if err != nil {
			fmt.Println("[ERROR] - Trying to open shm ", err)
		}
		//defer r.Close()
		
		// Try read the opened shared memory
		_, err = r.Read(rbuf)
		if err != nil {
			fmt.Println("[ERROR] - Trying to read shm ", err)
		}

		buf.Write(rbuf)
		err = binary.Read(buf, binary.LittleEndian, physics)
		if err != nil {
			fmt.Println("[ERROR] - Trying to write shm ", err)
		}

		err = r.Close()
		if err != nil {
			fmt.Println("[ERROR] - Trying to close open shared memory ", err)
		}
	}
}

func SendDataViaUdp() {
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
	fmt.Println("----- Running -----")

	physics := new(accshmdata.ACCPhysics)
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	// for {
	// 	<-ticker.C

	// 	// Reading will fail, if the game has not been started at least once
	// 	err := accshm.ReadPhysics(physics)
	// 	if err != nil {
	// 		fmt.Println("[ERROR] - Error reading frm client")
	// 	}

	wg.Add(1)
	go func() {
		<-ticker.C
		defer wg.Done()
		ReadSharedMemory(physics)
		fmt.Println("RPM: ", physics.RPM)
		fmt.Println("Gear: ", physics.Gear)
		fmt.Println("Brake: ", physics.Brake)
		fmt.Println("--------------------")
	}()

	// wg.Add(1)
	// go func() {
	// 	defer wg.Done()
	// 	SendDataViaUdp()
	// }()

	wg.Wait()
}