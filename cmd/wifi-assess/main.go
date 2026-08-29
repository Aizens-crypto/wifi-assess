package main

import (
	"fmt"
	"log"
	"os"
	"wifi-assess/internal/capture"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: wifi-assess <capture.pcapng>")
		os.Exit(1)
	}

	filename := os.Args[1]

	source := capture.NewPCAPSource(filename)

	if err := source.Open(); err != nil {
		log.Fatalf("failed to open capture: %v", err)
	}

	defer source.Close()

	packetCount := 0

	for {
		_, err := source.ReadPacket()

		if err != nil {
			break
		}

		packetCount++
	}

	fmt.Println("WiFi-Assess")
	fmt.Println("===========")
	fmt.Printf("Capture: %s\n", filename)
	fmt.Printf("Packets: %d\n", packetCount)
}
