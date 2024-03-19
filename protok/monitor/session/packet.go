package session

import (
	"fmt"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

type Packet struct {
	srcIP, dstIP                               string
	dstPort, srcPort                           uint16
	Seq, Ack                                   uint32
	Window                                     uint16
	FIN, SYN, RST, PSH, ACK, URG, ECE, CWR, NS bool
	payload                                    []byte
}

func NewPacket(gopkt gopacket.Packet) *Packet {
	// if gopkt.NetworkLayer() == nil {
	// 	return nil
	// }
	pkt := &Packet{}
	ipLayer := gopkt.Layer(layers.LayerTypeIPv4)
	if ipLayer != nil {
		ip, _ := ipLayer.(*layers.IPv4)
		pkt.dstIP = ip.DstIP.String()
		pkt.srcIP = ip.SrcIP.String()
	}
	tcpLayer := gopkt.Layer(layers.LayerTypeTCP)
	if tcpLayer != nil {
		tcp, _ := tcpLayer.(*layers.TCP)
		pkt.dstPort = uint16(tcp.DstPort)
		pkt.srcPort = uint16(tcp.SrcPort)
		pkt.Seq, pkt.Ack = tcp.Seq, tcp.Ack
		pkt.Window = tcp.Window
		pkt.FIN, pkt.SYN, pkt.RST, pkt.PSH, pkt.ACK, pkt.URG, pkt.ECE, pkt.CWR, pkt.NS = tcp.FIN, tcp.SYN, tcp.RST, tcp.PSH, tcp.ACK, tcp.URG, tcp.ECE, tcp.CWR, tcp.NS
		pkt.payload = tcp.Payload
	}
	return pkt
}

func (pkt *Packet) SrcAddr() string {
	return fmt.Sprintf("%s:%v", pkt.srcIP, pkt.srcPort)
}

func (pkt *Packet) SrcPort() uint16 {
	return pkt.srcPort
}

func (pkt *Packet) DstAddr() string {
	return fmt.Sprintf("%s:%v", pkt.dstIP, pkt.dstPort)
}

func (pkt *Packet) DstPort() uint16 {
	return pkt.dstPort
}

func (pkt *Packet) TransportAddr() string {
	return fmt.Sprintf("%s => %s", pkt.SrcAddr(), pkt.DstAddr())
}

func (pkt *Packet) Payload() []byte {
	return pkt.payload
}

func (pkt *Packet) IsBeginHandShake() bool {
	return pkt.SYN && !pkt.ACK
}

func (pkt Packet) IsBeginWaveHand() bool {
	return pkt.FIN && pkt.ACK
}

func (pkt Packet) IsReset() bool {
	return pkt.RST
}
