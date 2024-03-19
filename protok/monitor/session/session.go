package session

import (
	"fmt"

	"gitlab.kaiqitech.com/k7game/server/tools/goctl/protok/monitor/constants"
	"gitlab.kaiqitech.com/k7game/server/tools/goctl/protok/monitor/decoder"
	"gitlab.kaiqitech.com/k7game/server/tools/goctl/protok/monitor/decoder/k7proto"
)

type Session struct {
	cachePkts     []*Packet
	established   bool
	clientDecoder decoder.Decoder
	serverDecoder decoder.Decoder
	opts          Options
}

func NewSession(options ...Option) *Session {
	session := &Session{
		cachePkts: make([]*Packet, 0),
		opts:      NewOptions(options...),
	}
	switch session.opts.protocol {
	case "game", "conn":
		cliDecoder, err := k7proto.NewK7Decoder(session.opts.protocol)
		if err != nil {
			panic(err)
		}
		session.clientDecoder = cliDecoder
		svrDecoder, err := k7proto.NewK7Decoder(session.opts.protocol)
		if err != nil {
			panic(err)
		}
		session.serverDecoder = svrDecoder
	default:
		panic(constants.ErrInvalidDecoderType)
	}

	return session
}

func (s *Session) HandleNewTCPPacket(pkt *Packet) bool {
	payloadLen := len(pkt.Payload())
	// fmt.Printf("%s: window:%v, seq:%v, ack:%v, fin:%v, syn:%v, rst:%v, psh:%v, ack:%v, urg:%v, len:%v\n",
	// 	pkt.TransportAddr(), pkt.Window, pkt.Seq, pkt.Ack, pkt.FIN, pkt.SYN, pkt.RST, pkt.PSH, pkt.ACK, pkt.URG, payloadLen)
	if payloadLen == 0 {
		if pkt.IsBeginHandShake() {
			fmt.Printf("%s: 连接中...\n", pkt.TransportAddr())
		} else if pkt.IsBeginWaveHand() {
			fmt.Printf("%s: 断开连接...\n", pkt.TransportAddr())
			s.established = false
		} else if pkt.IsReset() {
			fmt.Printf("%s: 连接重置...\n", pkt.TransportAddr())
			s.established = false
		}
		if !s.established {
			s.cachePkts = append(s.cachePkts, pkt)
			return s.handleShakeAndWaveHand()
		}
		return true
	} else if !pkt.PSH && payloadLen == 1 {
		fmt.Printf("%s: window:%v, seq:%v, ack:%v, FIN:%v, SYN:%v, RST:%v, PSH:%v, ACK:%v, URG:%v, ECE:%v, CWR:%v, NS:%v, len:%v\n",
			pkt.TransportAddr(), pkt.Window, pkt.Seq, pkt.Ack, pkt.FIN, pkt.SYN, pkt.RST, pkt.PSH, pkt.ACK, pkt.URG, pkt.ECE, pkt.CWR, pkt.NS, payloadLen)
		return true
	}
	s.established = true
	var msg decoder.Packet
	var err error
	var decoder decoder.Decoder
	if s.opts.clientAddr == pkt.SrcAddr() {
		decoder = s.clientDecoder
	} else {
		decoder = s.serverDecoder
	}
	if pkt.PSH {
		msg, err = decoder.Decode(pkt.Payload())
	} else {
		decoder.Cache(pkt.Payload())
		return true
	}
	for {
		if err != nil {
			if err != constants.ErrPacketNotComplete {
				fmt.Printf("%s: decode failed: %s\n", pkt.TransportAddr(), err.Error())
			}
			break
		} else {
			isHeartbeat := msg.IsHeartbeat()
			if (isHeartbeat && s.opts.heartbeat) || (!isHeartbeat && (len(s.opts.filter) == 0 || msg.Filter(s.opts.filter))) {
				fmt.Printf("%s: %v\n", pkt.TransportAddr(), msg.ToString())
			}
		}
		msg, err = decoder.Decode([]byte{})
	}
	return true
}

func (s *Session) handleShakeAndWaveHand() bool {
	if len(s.cachePkts) >= 3 {
		if s.cachePkts[0].IsBeginHandShake() {
			if s.cachePkts[1].SYN && s.cachePkts[1].ACK && s.cachePkts[2].ACK {
				fmt.Printf("%s: 连接成功\n", s.cachePkts[2].TransportAddr())
				s.cachePkts = []*Packet{}
			}
		} else if s.cachePkts[0].IsBeginWaveHand() {
			s.cachePkts = []*Packet{}
			return false
		} else if s.cachePkts[0].IsReset() {
			s.cachePkts = []*Packet{}
			return false
		}
	}
	return true
}
