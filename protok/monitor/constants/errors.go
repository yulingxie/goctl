package constants

import "errors"

var (
	ErrDecryptFailed        = errors.New("decrypt failed")
	ErrInvalidHeaderSize    = errors.New("invalid header size")
	ErrPacketNotComplete    = errors.New("packet is not complete")
	ErrSocketIsAboutClose   = errors.New("socket is about to close")
	ErrInvalidDecoderType   = errors.New("invalid decoder type")
	ErrUnknownPacketVersion = errors.New("unknown packet version")
)
