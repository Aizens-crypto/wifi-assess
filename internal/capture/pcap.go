package capture

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

type PCAPSource struct {
	fileName     string
	handle       *pcap.Handle
	packetSource *gopacket.PacketSource
}

func NewPCAPSource(fileName string) *PCAPSource {

	return &PCAPSource{fileName: fileName}
}

func (p *PCAPSource) Open() error {

	handle, err := pcap.OpenOffline(p.fileName)

	if err != nil {
		return err
	}

	p.handle = handle
	linkType := handle.LinkType()

	p.packetSource = gopacket.NewPacketSource(handle, linkType)

	return nil

}

func (p *PCAPSource) ReadPacket() (gopacket.Packet, error) {
	return p.packetSource.NextPacket()
}

func (p *PCAPSource) Close() error {

	if p.handle != nil {

		p.handle.Close()
	}

	return nil
}
