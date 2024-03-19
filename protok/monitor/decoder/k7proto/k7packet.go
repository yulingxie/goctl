package k7proto

import (
	"fmt"
	"strconv"
	"strings"

	"gitlab.kaiqitech.com/nitro/nitro/v3/protok"
)

type K7Packet struct {
	option                uint32
	mainOpcode, subOpcode uint16
	registry              *protok.Registry
	k7Msg                 interface{}
}

func (pkt *K7Packet) Unmarshal(data []byte) error {
	pkt.k7Msg = pkt.registry.NewMsgByOpcode(pkt.mainOpcode, pkt.subOpcode)
	if pkt.k7Msg == nil {
		return fmt.Errorf("unknow opcode: %d %d", pkt.mainOpcode, pkt.subOpcode)
	}

	if err := protok.Unmarshal(data, pkt.k7Msg, (pkt.option&uint32(protok.PACKET_OPT_UTF8) != 0)); err != nil {
		return fmt.Errorf("pkt(%+v, %+v) Unmarshal err:%s", pkt.mainOpcode, pkt.subOpcode, err.Error())
	}
	return nil
}

func (pkt *K7Packet) ToString() string {
	return fmt.Sprintf("(%4d,%4d): %s ", pkt.mainOpcode, pkt.subOpcode, protok.ToString(pkt.k7Msg, nil))
}

func (pkt *K7Packet) IsHeartbeat() bool {
	return pkt.mainOpcode == 0 && (pkt.subOpcode == 1 || pkt.subOpcode == 5)
}

func (pkt *K7Packet) Filter(filter []string) bool {
	for _, v := range filter {
		opcodes := strings.Split(v, ",")
		if len(opcodes) < 2 {
			continue
		}
		if opcodes[0] == strconv.Itoa(int(pkt.mainOpcode)) && opcodes[1] == strconv.Itoa(int(pkt.subOpcode)) {
			return true
		}
	}
	return false
}

func NewK7Packet(mainCode, subCode uint16, option uint32, registry *protok.Registry) *K7Packet {
	return &K7Packet{
		option:     option,
		mainOpcode: mainCode,
		subOpcode:  subCode,
		registry:   registry,
	}
}
