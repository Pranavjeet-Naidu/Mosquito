package steg

import (
    "bytes"
    "encoding/binary"
)

const (
    // MagicByte identifies a Mosquito steganography header
    MagicByte byte = 0x53
    // Version of the header format
    Version byte = 0x02
)

// MessageFlags for different payload types and features
type MessageFlags byte

const (
    // FlagEncrypted indicates the payload is encrypted
    FlagEncrypted MessageFlags = 1 << iota
    // FlagCompressed indicates the payload is compressed
    FlagCompressed
    // FlagImage indicates the payload is an image
    FlagImage
)

// Header represents the metadata for a hidden message
type Header struct {
    Magic     byte        // Magic byte (0x53)
    Version   byte        // Header version
    Mode      StegMode    // Steganography mode used
    Flags     MessageFlags // Payload flags
    PayloadLen uint32      // Length of the payload in bytes
}

// MarshalHeader converts a header to bytes
func MarshalHeader(h Header) []byte {
    buf := new(bytes.Buffer)
    buf.WriteByte(h.Magic)
    buf.WriteByte(h.Version)
    buf.WriteByte(byte(h.Mode))
    buf.WriteByte(byte(h.Flags))
    binary.Write(buf, binary.BigEndian, h.PayloadLen)
    return buf.Bytes()
}

// UnmarshalHeader parses bytes into a header
func UnmarshalHeader(data []byte) (Header, error) {
    if len(data) < 8 {
        return Header{}, ErrInvalidHeader
    }

    h := Header{
        Magic:     data[0],
        Version:   data[1],
        Mode:      StegMode(data[2]),
        Flags:     MessageFlags(data[3]),
        PayloadLen: binary.BigEndian.Uint32(data[4:8]),
    }

    if h.Magic != MagicByte {
        return Header{}, ErrInvalidMagic
    }

    return h, nil
}

// Size returns the size of the header in bytes
func (h Header) Size() int {
    return 8 // Magic(1) + Version(1) + Mode(1) + Flags(1) + PayloadLen(4)
}

// IsEncrypted returns true if the payload is encrypted
func (h Header) IsEncrypted() bool {
    return (h.Flags & FlagEncrypted) != 0
}

// IsCompressed returns true if the payload is compressed
func (h Header) IsCompressed() bool {
    return (h.Flags & FlagCompressed) != 0
}

// IsImage returns true if the payload is an image
func (h Header) IsImage() bool {
    return (h.Flags & FlagImage) != 0
}