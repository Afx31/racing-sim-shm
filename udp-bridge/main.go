package main

import (
	"encoding/binary"
	"fmt"
	"math"
	"net"
	"sync"
	"unsafe"

	"udp-bridge/internal/accshmdata"

	"github.com/hidez8891/shm"
)

const (
	shmNamePhysics    = "Local\\acpmf_physics"
	shmNameGraphics		= "Local\\acpmf_graphics"
	shmNameStatic			= "Local\\acpmf_static"
	SERVER_UDP_ADDR 	= ""
	
	PACKET_LENGTH			= 100
)

var (
	wg sync.WaitGroup
)

func ReadSharedMemory(physics *accshmdata.ACCPhysics) {
	fmt.Println("--- Reading Memory ---")
	
	for {
		physicsSize := (int32)(unsafe.Sizeof(*physics))
		
		// Open shared memory (shm reader)
		r, err := shm.Open(shmNamePhysics, physicsSize)
		if err != nil {
			fmt.Println("[ERROR] - Trying to open shm ", err)
		}
		defer r.Close()
		

		//----------------------------------------------------------------
		// TODO: Currently testing Solution 1 & 2 below

		// ---- Solution 1
		// Read the opened shared memory into a raw buffer
		// rbuf := make([]byte, physicsSize)
		// _, err = r.Read(rbuf)
		// if err != nil {
		// 	fmt.Println("[ERROR] - Trying to read shm ", err)
		// }

		// Convert into a buffer wrapper for to decode into the struct
		// buf := &bytes.Buffer{}
		// buf.Write(rbuf)
		// err = binary.Read(buf, binary.LittleEndian, physics)
		// if err != nil {
		// 	fmt.Println("[ERROR] - Trying to write shm ", err)
		// }

		// ---- Solution 2
		err = binary.Read(r, binary.LittleEndian, physics)
		if err != nil {
			fmt.Println("[ERROR] - Trying to write shm ", err)
		}

		// ----------------------------------------------------------------
		

		err = r.Close()
		if err != nil {
			fmt.Println("[ERROR] - Trying to close open shared memory ", err)
		}
	}
}

func SendDataViaUdp(physics *accshmdata.ACCPhysics) {
	conn, err := net.Dial("udp", SERVER_UDP_ADDR)
	if err != nil {
		fmt.Println("[ERROR] - Connecting to UDP Server: ", err)
	}
	defer conn.Close()

	counter := 0

	for {
		b := make([]byte, PACKET_LENGTH)

		binary.LittleEndian.PutUint32(b[0:4], uint32(physics.PacketId))
		binary.LittleEndian.PutUint32(b[4:8], uint32(physics.RPM))
		binary.LittleEndian.PutUint32(b[8:12], math.Float32bits(physics.SpeedKmh))
		binary.LittleEndian.PutUint32(b[12:16], uint32(physics.Gear))
		binary.LittleEndian.PutUint32(b[16:20], math.Float32bits(physics.Gas * 100))
		binary.LittleEndian.PutUint32(b[20:24], math.Float32bits(physics.Brake * 100))

		_, err = conn.Write(b)
		if err != nil {
			fmt.Println("[ERROR] - Writing to UDP Server: ", err)
		}
    
		counter++
	}
}

func main() {
	fmt.Println("----- Running -----")

	physics := new(accshmdata.ACCPhysics)

	wg.Add(1)
	go func() {
		defer wg.Done()
		ReadSharedMemory(physics)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		SendDataViaUdp(physics)
	}()

	wg.Wait()
}