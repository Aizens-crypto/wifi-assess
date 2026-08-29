package capture

import (
	"github.com/google/gopacket"
)

type Source interface {
	Open() error
	ReadPacket() (gopacket.Packet, error)
	Close() error
}
