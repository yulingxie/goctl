package monitor

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"gitlab.kaiqitech.com/k7game/server/tools/goctl/protok/monitor/session"
	"gitlab.kaiqitech.com/k7game/server/tools/goctl/protok/monitor/ui"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

type Monitor struct {
	opts     Options
	sessions map[string]*session.Session
}

func NewMonitor(options ...Option) *Monitor {
	return &Monitor{
		opts:     NewOptions(options...),
		sessions: make(map[string]*session.Session),
	}
}

func (m Monitor) SelectDevice() string {
	devices, err := pcap.FindAllDevs()
	if err != nil {
		panic(err)
	}
	metaData := map[string]interface{}{}
	for _, device := range devices {
		for _, v := range device.Addresses {
			ip := v.IP.To4()
			if ip != nil {
				key := fmt.Sprintf("%15s\t(%s)", ip.String(), device.Description)
				metaData[key] = device.Name
			}
		}
	}
	if selected := ui.NewListModel("选择监听IP:", metaData).Select(); selected != nil {
		return selected.(string)
	}
	return ""
}

func (m *Monitor) MonitorDevice(deviceName string) {
	fmt.Println("Monitor starting ...")
	handle, err := pcap.OpenLive(deviceName, m.opts.snapshotLen, m.opts.promiscuous, m.opts.timeout)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()
	exitCh := make(chan os.Signal, 1)
	signal.Notify(exitCh, syscall.SIGINT, syscall.SIGTERM)
	m.ServeTCPacket(handle, exitCh)

	fmt.Println("Monitor end")
}

func (m *Monitor) ServeTCPacket(handle *pcap.Handle, exitSig chan os.Signal) {
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for {
		select {
		case pkt := <-packetSource.Packets():
			m.handleNewPacket(pkt)
		case <-exitSig:
			return
		}
	}
}

func (m *Monitor) handleNewPacket(gopkt gopacket.Packet) bool {
	pkt := session.NewPacket(gopkt)
	if pkt == nil {
		return false
	}
	if pkt.DstAddr() != m.opts.serverAddr && pkt.SrcAddr() != m.opts.serverAddr {
		return false
	}
	clientAddr := pkt.SrcAddr()
	if pkt.DstAddr() != m.opts.serverAddr {
		clientAddr = pkt.DstAddr()
	}
	if _, ok := m.sessions[clientAddr]; !ok {
		m.sessions[clientAddr] = session.NewSession(
			session.WithProtocol(m.opts.protocol),
			session.WithClientAddr(clientAddr),
			session.WithServerAddr(m.opts.serverAddr),
			session.WithFilter(m.opts.packets),
			session.WithHeartbeat(m.opts.heartbeat),
		)
	}
	if !m.sessions[clientAddr].HandleNewTCPPacket(pkt) {
		delete(m.sessions, clientAddr)
	}
	return true
}

func (m *Monitor) dumpPacketInfo(gopkt gopacket.Packet) {
	// Let's see if the packet is an ethernet packet
	// 判断数据包是否为以太网数据包，可解析出源mac地址、目的mac地址、以太网类型（如ip类型）等
	ethernetLayer := gopkt.Layer(layers.LayerTypeEthernet)
	if ethernetLayer != nil {
		fmt.Println("Ethernet layer detected.")
		ethernetPacket, _ := ethernetLayer.(*layers.Ethernet)
		fmt.Println("Source MAC: ", ethernetPacket.SrcMAC)
		fmt.Println("Destination MAC: ", ethernetPacket.DstMAC)
		// Ethernet type is typically IPv4 but could be ARP or other
		fmt.Println("Ethernet type: ", ethernetPacket.EthernetType)
		fmt.Println()
	}
	// Let's see if the packet is IP (even though the ether type told us)
	// 判断数据包是否为IP数据包，可解析出源ip、目的ip、协议号等
	ipLayer := gopkt.Layer(layers.LayerTypeIPv4)
	if ipLayer != nil {
		fmt.Println("IPv4 layer detected.")
		ip, _ := ipLayer.(*layers.IPv4)
		// IP layer variables:
		// Version (Either 4 or 6)
		// IHL (IP Header Length in 32-bit words)
		// TOS, Length, Id, Flags, FragOffset, TTL, Protocol (TCP?),
		// Checksum, SrcIP, DstIP
		fmt.Printf("From %s to %s\n", ip.SrcIP, ip.DstIP)
		fmt.Println("Protocol: ", ip.Protocol)
		fmt.Println()
	}
	// Let's see if the packet is TCP
	// 判断数据包是否为TCP数据包，可解析源端口、目的端口、seq序列号、tcp标志位等
	tcpLayer := gopkt.Layer(layers.LayerTypeTCP)
	if tcpLayer != nil {
		fmt.Println("TCP layer detected.")
		tcp, _ := tcpLayer.(*layers.TCP)
		// TCP layer variables:
		// SrcPort, DstPort, Seq, Ack, DataOffset, Window, Checksum, Urgent
		// Bool flags: FIN, SYN, RST, PSH, ACK, URG, ECE, CWR, NS
		fmt.Printf("From port %d to %d\n", tcp.SrcPort, tcp.DstPort)
		fmt.Println("Sequence number: ", tcp.Seq)
		fmt.Println()
	}
	// Iterate over all layers, printing out each layer type
	fmt.Println("All packet layers:")
	for _, layer := range gopkt.Layers() {
		fmt.Println("- ", layer.LayerType())
	}
	///.......................................................
	// Check for errors
	// 判断layer是否存在错误
	if err := gopkt.ErrorLayer(); err != nil {
		fmt.Println("Error decoding some part of the packet:", err)
	}
}
