package decoder

type Packet interface {
	ToString() string
	Unmarshal([]byte) error
	Filter([]string) bool
	IsHeartbeat() bool
}

type Decoder interface {
	Cache([]byte)
	Decode([]byte) (Packet, error)
}
