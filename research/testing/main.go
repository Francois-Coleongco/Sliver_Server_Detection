package main

import (
	"fmt"
	"github.com/gopacket/gopacket"
	"github.com/gopacket/gopacket/pcap"
	"log"
	"os"
)

func create_and_setup_logs() {
	file, err := os.OpenFile("sniffy.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}

	log.SetOutput(file)
}

func main() {
	// Open the pcap file or live capture

	create_and_setup_logs()

	// dest_IP := retrieve_my_ip()

	filter := fmt.Sprintf("tcp port 8888") // sliver listens on port 8888 by default. just hardcoding this here for now just for testing myself. it is possible to change the port i am aware, just rn for testing i wanna make it hardcoded

	fmt.Println(filter)

	// Create a packet, but don't actually decode anything yet
	if handle, err := pcap.OpenLive("enp0s3", 1600, false, pcap.BlockForever); err != nil {
		panic(err)
	} else if err := handle.SetBPFFilter(filter); err != nil { // optional

		panic(err)
	} else {
		packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
		for packet := range packetSource.Packets() {
			fmt.Println(packet.String())

			log.Println(packet)

		}
	}
	// Now, decode the packet up to the first IPv4 layer found but no further.
	// If no IPv4 layer was found, the whole packet will be decoded looking for
	// it.
}
