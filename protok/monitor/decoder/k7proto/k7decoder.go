package k7proto

import (
	"encoding/binary"
	"fmt"
	"strings"

	"gitlab.kaiqitech.com/k7game/server/tools/goctl/protok/monitor/decoder"

	"gitlab.kaiqitech.com/k7game/server/tools/goctl/protok/monitor/constants"

	"gitlab.kaiqitech.com/k7game/server/components/seraph.git/server/connsvr"
	"gitlab.kaiqitech.com/k7game/server/components/seraph.git/server/gamesvr"
	"gitlab.kaiqitech.com/nitro/nitro/v3/protok"
)

var (
	minPacketSize      = protok.PACKET_HEADER_SIZE + protok.PACKET_TAIL_SIZE
	websocketClientSig = []uint8{0x47, 0x45, 0x54, 0x20, 0x2f, 0x20, 0x48, 0x54, 0x54, 0x50, 0x2f, 0x31, 0x2e, 0x31, 0x0d, 0x0a}
	websocketServerSig = []uint8{0x48, 0x54, 0x54, 0x50, 0x2f, 0x31, 0x2e, 0x31, 0x20, 0x31, 0x30, 0x31, 0x20, 0x57, 0x65, 0x62, 0x53, 0x6f, 0x63, 0x6b, 0x65, 0x74, 0x20, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x20, 0x48, 0x61, 0x6e, 0x64, 0x73, 0x68, 0x61, 0x6b, 0x65, 0xd, 0xa}
)

type K7Decoder struct {
	registry       *protok.Registry
	crypt          protok.Crypt
	cacheData      []byte
	pktEncaseType  int
	headerDetected bool
}

func NewK7Decoder(proto string) (*K7Decoder, error) {
	d := &K7Decoder{
		cacheData:      make([]byte, 0),
		pktEncaseType:  constants.PKT_ENCASE_T_NONE,
		headerDetected: false,
	}
	switch proto {
	case "conn":
		d.registry = protok.NewRegistry(connsvr.ConnProtos)
	case "game":
		d.registry = protok.NewRegistry(gamesvr.GameProtos)
	default:
		return nil, constants.ErrInvalidDecoderType
	}
	return d, nil
}

func (d *K7Decoder) Cache(data []byte) {
	d.cacheData = append(d.cacheData, data...)
}

func (d *K7Decoder) Decode(data []byte) (decoder.Packet, error) {
	d.cacheData = append(d.cacheData, data...)
	return d.decode()
}

func (d *K7Decoder) detectHeader() error {
	if d.headerDetected {
		return nil
	}
	return d.detectWebSocketHandshake()
}

func (d *K7Decoder) detectWebSocketHandshake() error {
	detectSig := func(sig []uint8) (bool, error) {
		if len(d.cacheData) < len(sig) {
			return false, constants.ErrPacketNotComplete
		}
		if string(d.cacheData[:len(sig)]) != string(sig) {
			return false, nil
		}
		return true, nil
	}
	if detect, err := detectSig(websocketClientSig); !detect && err == nil {
		if detect, err = detectSig(websocketServerSig); err != nil {
			return err
		} else if !detect && err == nil {
			return nil
		}
	} else if err != nil {
		return err
	}
	d.headerDetected = true
	d.pktEncaseType = constants.PKT_ENCASE_T_WEBSOCKET
	idx := strings.Index(string(d.cacheData), "\r\n\r\n")
	if idx == -1 {
		return constants.ErrPacketNotComplete
	}
	d.cacheData = d.cacheData[idx+4:]
	d.headerDetected = true
	return nil
}

func (d *K7Decoder) detectWebSocket() error {
	if len(d.cacheData) < 2 {
		return constants.ErrPacketNotComplete
	}
	var fin, hasMask bool
	headerSize := 2
	var headerPayloadLen uint16
	{
		fin = (d.cacheData[0] & 0x80) != 0
		opcode := d.cacheData[0] & 0x0f
		byteData := d.cacheData[1]
		hasMask = (byteData & 0x80) != 0
		headerPayloadLen = uint16(byteData & 0x7f)
		if opcode == 0x01 || opcode == 0x08 {
			return constants.ErrSocketIsAboutClose
		}
		if hasMask {
			headerSize += 4
		}
		if headerPayloadLen == 126 {
			headerSize += 2
		} else if headerPayloadLen == 127 {
			headerSize += 8
		}
	}
	if len(d.cacheData) < headerSize {
		return constants.ErrPacketNotComplete
	}
	payloadLen := uint64(headerPayloadLen)
	if headerPayloadLen == 126 {
		payloadLen = uint64(binary.BigEndian.Uint16(d.cacheData[2:]))
	} else if headerPayloadLen == 127 {
		payloadLen = binary.BigEndian.Uint64(d.cacheData[2:])
	}

	data := d.cacheData[headerSize:]
	dataLen := len(data)
	if hasMask {
		maskOffset := 2
		if payloadLen == 126 {
			maskOffset += 2
		} else if payloadLen == 127 {
			maskOffset += 8
		}
		mask := d.cacheData[maskOffset : maskOffset+4]
		if dataLen > 0 {
			for i := 0; i < dataLen; i++ {
				data[i] = data[i] ^ mask[i%4]
			}
		}
	}

	if fin && payloadLen != 0 {
		if dataLen < 2 {
			return constants.ErrSocketIsAboutClose
		}
		d.cacheData = data
	}
	return nil
}

func (d *K7Decoder) checkPacket() (decoder.Packet, error) {
	if len(d.cacheData) < 2 {
		return nil, constants.ErrPacketNotComplete
	}
	headerSize := binary.LittleEndian.Uint16(d.cacheData)
	if headerSize != protok.PACKET_HEADER_SIZE {
		return nil, constants.ErrInvalidHeaderSize
	}
	if len(d.cacheData) < int(minPacketSize) {
		// 不完整 min size
		return nil, constants.ErrPacketNotComplete
	}
	tmpData := []byte{}
	tmpData = append(tmpData, d.cacheData...)
	pktVersion := binary.LittleEndian.Uint32(tmpData[4:])
	if pktVersion != protok.PACKET_VERSION {
		return nil, constants.ErrUnknownPacketVersion
	}
	pkgOption := binary.LittleEndian.Uint32(tmpData[8:])
	// 整包大小
	pkgSize := binary.LittleEndian.Uint16(tmpData[12:])
	cmdLen := pkgSize - minPacketSize
	if len(tmpData) < int(pkgSize) {
		// 不完整 full size
		return nil, constants.ErrPacketNotComplete
	}
	d.cacheData = make([]byte, 0)
	cmdData := tmpData[protok.PACKET_HEADER_SIZE : protok.PACKET_HEADER_SIZE+cmdLen]
	pktLen := len(cmdData) + int(minPacketSize)
	if len(tmpData) > pktLen {
		// 数据多
		d.cacheData = append(d.cacheData, tmpData[pktLen:]...)
	}
	_, decryptData := d.crypt.DecryptRecv(cmdData)
	mainOpcode := binary.LittleEndian.Uint16(decryptData[8:])
	subOpcode := binary.LittleEndian.Uint16(decryptData[10:])
	if (pkgOption & uint32(protok.PACKET_OPT_ZIP)) != 0 {
		var err error
		decryptData, err = protok.ZipDecompress(decryptData[20:])
		if err != nil {
			return nil, fmt.Errorf("pkt(%+v, %+v) ZipDecompress err:%s", mainOpcode, subOpcode, err.Error())
		}
	} else {
		decryptData = decryptData[16:]
	}
	pkt := NewK7Packet(mainOpcode, subOpcode, pkgOption, d.registry)
	if err := pkt.Unmarshal(decryptData); err != nil {
		return nil, err
	}
	return pkt, nil
}

func (d *K7Decoder) decode() (decoder.Packet, error) {
	if err := d.detectHeader(); err != nil {
		return nil, err
	}
	switch d.pktEncaseType {
	case constants.PKT_ENCASE_T_WEBSOCKET:
		if err := d.detectWebSocket(); err != nil {
			return nil, err
		}
	}
	return d.checkPacket()
}
