package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"sort"

	"wifi-assess/internal/capture"
	"wifi-assess/internal/discovery"
	"wifi-assess/internal/wifi"
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
	skipped := 0
	frameCounts := make(map[string]int)
	tracker := discovery.NewTracker()

	for {
		pkt, err := source.ReadPacket()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("read error: %v", err)
			continue
		}

		packetCount++

		parsed, err := wifi.ParsePacket(pkt)
		if err != nil {
			// Not every captured packet is 802.11 (e.g. loopback/junk
			// frames can end up in a capture). Count and move on.
			skipped++
			continue
		}

		frameCounts[parsed.FrameType]++
		tracker.Observe(parsed)
	}

	fmt.Println("WiFi-Assess")
	fmt.Println("===========")
	fmt.Printf("Capture: %s\n", filename)
	fmt.Printf("Packets: %d\n", packetCount)
	fmt.Printf("Non-802.11 / unparsed: %d\n", skipped)
	fmt.Println()
	fmt.Println("Frame types:")

	types := make([]string, 0, len(frameCounts))
	for t := range frameCounts {
		types = append(types, t)
	}
	sort.Strings(types)

	for _, t := range types {
		fmt.Printf("  %-30s %d\n", t, frameCounts[t])
	}

	aps := tracker.AccessPoints()
	sort.Slice(aps, func(i, j int) bool { return aps[i].BSSID < aps[j].BSSID })

	fmt.Println()
	fmt.Printf("Access points (%d):\n", len(aps))
	for _, ap := range aps {
		privacy := "unknown"
		if ap.HasCapabilityInfo {
			if ap.PrivacyEnabled {
				privacy = "encrypted"
			} else {
				privacy = "open"
			}
		}
		ssid := ap.SSID
		if ssid == "" {
			ssid = "<hidden/unknown>"
		}
		fmt.Printf("  %-17s  ch=%-3d  %-9s  beacons=%-5d probe_resp=%-5d  %s\n",
			ap.BSSID, ap.Channel, privacy, ap.BeaconCount, ap.ProbeResponseCount, ssid)
	}

	clients := tracker.Clients()
	sort.Slice(clients, func(i, j int) bool { return clients[i].MAC < clients[j].MAC })

	fmt.Println()
	fmt.Printf("Clients (%d):\n", len(clients))
	for _, c := range clients {
		assoc := c.AssociatedBSSID
		if assoc == "" {
			assoc = "-"
		}
		fmt.Printf("  %-17s  probes=%-5d  assoc=%s  probed_ssids=%v\n",
			c.MAC, c.ProbeRequestCount, assoc, c.ProbedSSIDs)
	}
}
