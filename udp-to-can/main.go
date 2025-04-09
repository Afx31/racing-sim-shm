package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"net"
)

const udpAddress = ""

func main() {
	fmt.Println("--- Running.. ---")
	
	conn, err := net.ListenPacket("udp", udpAddress)
	if err != nil {
			log.Fatal("Error listening on UDP port:", err)
	}
	defer conn.Close()

	b := make([]byte, 30)
	counter := 0
	for {
		_, _, err := conn.ReadFrom(b)
		if err != nil {
				log.Fatal("Error reading from UDP:", err)
		}

		packetId := binary.LittleEndian.Uint32(b[0:4])
		rpm := binary.LittleEndian.Uint32(b[4:8])
		speed := math.Float32frombits(binary.LittleEndian.Uint32(b[8:12]))
		gear := binary.LittleEndian.Uint32(b[12:16])
		tps := math.Float32frombits(binary.LittleEndian.Uint32(b[16:20]))

		fmt.Println("-------------------------------------------------")
		fmt.Println(counter)
		fmt.Println(b)
		fmt.Println("PacketId: ", packetId)
		fmt.Println("RPM: ", rpm)
		fmt.Println("SpeedKmh: ", speed)
		fmt.Println("Gear: ", gear)
		fmt.Println("Gas: ", tps)
		
		counter++
	}
}
