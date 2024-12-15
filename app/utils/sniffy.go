package utils

import (
	"bytes"
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

func Sniffer(my_port string, pid string, pid_chan chan string) {

	fmt.Println("EYYY WE GOT HERE")

	create_and_setup_logs()
	// dest_IP := retrieve_my_ip()

	filter := fmt.Sprintf("tcp dst port %s", my_port) // sliver listens on port 8888 by default. just hardcoding this here for now just for testing myself. it is possible to change the port i am aware, just rn for testing i wanna make it hardcoded

	fmt.Println(filter)

	// Create a packet, but don't actually decode anything yet
	if handle, err := pcap.OpenLive("wlo1", 1500, false, pcap.BlockForever); err != nil {
		panic(err)
	} else if err := handle.SetBPFFilter(filter); err != nil { // optional

		panic(err)
	} else {

		app_data_bytes := []byte{23, 3, 3} // remember this is base 10 now
		packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
		for packet := range packetSource.Packets() {

			data := packet.TransportLayer().LayerPayload()
			log.Println(packet, data)

			if bytes.Contains(data, app_data_bytes) {
				// contains application data

				fmt.Println("FOUND APP_DATA")

				pid_chan <- pid // threshold reached
			}

		}
	}
	// Now, decode the packet up to the first IPv4 layer found but no further.
	// If no IPv4 layer was found, the whole packet will be decoded looking for
	// it.
	pid_chan <- "" // only way for this to happen would be it couldn't communicate with the interface or something. we want to keep the sniffer running on the processes until we know they communicate under encryption
}
