package main

import (
	"context"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"net"
	"sync"

	// "time"

	"udp-to-can/internal/hondata"

	"go.einride.tech/can"
	"go.einride.tech/can/pkg/socketcan"
)

const (
	udpAddress = ""
	SETTINGS_TICKER = 100
)

var (
	wg sync.WaitGroup
	frame660 = hondata.Frame660{}
	frame662 = hondata.Frame662{}
)

func ReadDataFromUdp() {
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

		// packetId := binary.LittleEndian.Uint32(b[0:4])
		rpm := binary.LittleEndian.Uint32(b[4:8])
		speed := math.Float32frombits(binary.LittleEndian.Uint32(b[8:12]))
		gear := binary.LittleEndian.Uint32(b[12:16])
		tps := math.Float32frombits(binary.LittleEndian.Uint32(b[16:20]))

		frame660.Rpm = uint16(rpm)
		frame660.Speed = uint16(speed)
		frame660.Gear = uint8(gear)

		frame662.Tps = uint16(tps)
		
		counter++
	}
}

func SendDataToCan() {
	conn, err := socketcan.DialContext(context.Background(), "can", "vcan0")
	if err != nil {
		fmt.Println("[ERROR] - Cannot connect to vcan0")
	}
	defer conn.Close()

	tx := socketcan.NewTransmitter(conn)

	// ticker := time.NewTicker(SETTINGS_TICKER * time.Millisecond)
	// defer ticker.Stop()

	counter := 0

	for {
		// <-ticker.C
		switch (counter) {
		case 0:
			f660 := can.Frame {
				ID: 660,
				Length: 8,
				Data: [8]byte { 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00 },
			}
			binary.BigEndian.PutUint16(f660.Data[0:2], frame660.Rpm)
			binary.BigEndian.PutUint16(f660.Data[2:4], frame660.Speed)
			f660.Data[4] = frame660.Gear

			_ = tx.TransmitFrame(context.Background(), f660)
				// fmt.Println("Sent 660: ", f660)

		case 1:
			f662 := can.Frame {
				ID: 662,
				Length: 8,
				Data: [8]byte { 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00 },
			}
			binary.BigEndian.PutUint16(f662.Data[0:2], frame662.Tps)

			_ = tx.TransmitFrame(context.Background(), f662)
				// fmt.Println("Sent 662: ", f662)
		}
		
		if (counter == 2) {
			counter = 0
		} else {
			counter++
		}
	}

}

func main() {
	fmt.Println("--- Running ---")
	
	wg.Add(1)
	go func() {
		defer wg.Done()
		ReadDataFromUdp()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		SendDataToCan()
	}()

	wg.Wait()
}
